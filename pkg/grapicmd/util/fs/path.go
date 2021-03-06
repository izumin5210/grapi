package fs

import (
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

type getOSUserFunc func() (*user.User, error)

const (
	// PackageSeparator is a package separator string on protobuf.
	PackageSeparator = "."
)

// Make visible for testing
var (
	BuildContext build.Context
	GetOSUser    getOSUserFunc
)

func init() {
	BuildContext = build.Default
	GetOSUser = user.Current
}

// GetImportPath creates the golang package path from the given path.
func GetImportPath(rootPath string) (importPath string, err error) {
	for _, gopath := range filepath.SplitList(BuildContext.GOPATH) {
		prefix := filepath.Join(gopath, "src") + string(filepath.Separator)
		// FIXME: should not use strings.HasPrefix
		if strings.HasPrefix(rootPath, prefix) {
			importPath = filepath.ToSlash(strings.Replace(rootPath, prefix, "", 1))
			break
		}
	}
	if importPath == "" {
		importPath = filepath.Base(rootPath)
	}
	return
}

// GetPackageName generates the package name of this application from the given path and envs.
func GetPackageName(rootPath string) (string, error) {
	importPath, err := GetImportPath(rootPath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	entries := strings.Split(importPath, string(filepath.Separator))
	if len(entries) < 2 {
		u, err := GetOSUser()
		if err != nil {
			return "", errors.WithStack(err)
		}
		entries = []string{u.Username, entries[0]}
	}
	entries = entries[len(entries)-2:]
	if strings.Contains(entries[0], PackageSeparator) {
		s := strings.Split(entries[0], PackageSeparator)
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
		entries[0] = strings.Join(s, PackageSeparator)
	}
	pkgName := strings.Join(entries[len(entries)-2:], PackageSeparator)
	pkgName = strings.Replace(pkgName, "-", "_", -1)
	return pkgName, nil
}

// FindMainPackagesAndSources returns go source file names by main package directories.
func FindMainPackagesAndSources(fs afero.Fs, dir string) (map[string][]string, error) {
	out := make(map[string][]string)
	fset := token.NewFileSet()
	err := afero.Walk(fs, dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		if info.IsDir() || filepath.Ext(info.Name()) != ".go" || strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}
		data, err := afero.ReadFile(fs, path)
		if err != nil {
			zap.L().Warn("failed to read a file", zap.Error(err), zap.String("path", path))
			return nil
		}
		f, err := parser.ParseFile(fset, "", data, parser.PackageClauseOnly)
		if err != nil {
			zap.L().Warn("failed to parse a file", zap.Error(err), zap.String("path", path), zap.String("body", string(data)))
			return nil
		}
		if f.Package.IsValid() && f.Name.Name == "main" {
			dir := filepath.Dir(path)
			out[dir] = append(out[dir], info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return out, nil
}
