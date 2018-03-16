package fs

import (
	"go/build"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
)

var (
	// BuildContext is a build context object.
	BuildContext build.Context
)

func init() {
	BuildContext = build.Default
}

// GetImportPath creates the golang package path from the given path.
func GetImportPath(rootPath string) (importPath string, err error) {
	for _, gopath := range filepath.SplitList(BuildContext.GOPATH) {
		prefix := filepath.Join(gopath, "src") + "/"
		// FIXME: should not use strings.HasPrefix
		if strings.HasPrefix(rootPath, prefix) {
			importPath = strings.Replace(rootPath, prefix, "", 1)
			break
		}
	}
	if importPath == "" {
		err = errors.New("failed to get the import path")
	}
	return
}

var (
	requiredFiles = []string{"grapi.toml"}
	requiredDirs  = []string{"app", "api", "cmd"}
)

// LookupRoot returns the application's root directory if the current directory is inside of a grapi application.
func LookupRoot(fs afero.Fs, dir string) (string, bool) {
	var eg errgroup.Group

	for _, f := range requiredFiles {
		f := f
		eg.Go(func() error {
			ok, err := afero.Exists(fs, filepath.Join(dir, f))
			if err != nil || !ok {
				return errors.Errorf("%s does not exist", f)
			}
			return nil
		})
	}

	for _, d := range requiredDirs {
		d := d
		eg.Go(func() error {
			ok, err := afero.DirExists(fs, filepath.Join(dir, d))
			if err != nil || !ok {
				return errors.Errorf("%s does not exist", d)
			}
			return nil
		})
	}

	if eg.Wait() == nil {
		return dir, true
	}

	if dir == "/" {
		return "", false
	}

	return LookupRoot(fs, filepath.Dir(dir))
}

// FindMainPackagesAndSources returns go source file names by main package directories.
func FindMainPackagesAndSources(fs afero.Fs, dir string) (map[string][]string, error) {
	out := make(map[string][]string)
	return out, nil
}
