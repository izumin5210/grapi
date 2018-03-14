package usecase

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/serenize/snaker"

	"github.com/izumin5210/clicontrib/clog"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// GenerateServiceUsecase is an useecase interface for geenrating .proto file and its implementation skeleton.
type GenerateServiceUsecase interface {
	Generate(path string) error
	Destroy(path string) error
}

type generateServiceUsecase struct {
	ui        module.UI
	generator module.Generator
	rootDir   string
}

// NewGenerateServiceUsecase returns an new GenerateServiceUsecase implementation instance.
func NewGenerateServiceUsecase(ui module.UI, generator module.Generator, rootDir string) GenerateServiceUsecase {
	return &generateServiceUsecase{
		ui:        ui,
		generator: generator,
		rootDir:   rootDir,
	}
}

func (u *generateServiceUsecase) Generate(path string) error {
	data, err := u.createParams(path)
	if err != nil {
		return errors.WithStack(err)
	}

	u.ui.Section("Generate service")
	clog.Debug("Generate service", "params", data)
	return u.generator.Generate(u.rootDir, data)
}

func (u *generateServiceUsecase) Destroy(path string) error {
	data, err := u.createParams(path)
	if err != nil {
		return errors.WithStack(err)
	}

	u.ui.Section("Destroy service")
	clog.Debug("Destroy service", "params", data)
	return u.generator.Destroy(u.rootDir, data)
}

func (u *generateServiceUsecase) createParams(path string) (map[string]interface{}, error) {
	// github.com/foo/bar
	importPath, err := fs.GetImportPath(u.rootDir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// path => baz/qux/quux
	path = strings.Replace(path, "-", "_", -1)

	// quux
	name := filepath.Base(path)
	// Quux
	serviceName := snaker.SnakeToCamel(name)

	// baz/qux
	packagePath := filepath.Dir(path)
	// qux
	packageName := filepath.Base(packagePath)

	// api/baz/qux
	pbgoPackagePath := filepath.Join("api", packagePath)
	// qux_pb
	pbgoPackageName := filepath.Base(pbgoPackagePath) + "_pb"

	if packagePath == "." {
		packagePath = "server"
		packageName = packagePath
		pbgoPackagePath = "api"
		pbgoPackageName = pbgoPackagePath + "_pb"
	}

	protoPackageChunks := []string{}
	for _, pkg := range strings.Split(filepath.Join(importPath, "api", filepath.Dir(path)), "/") {
		chunks := strings.Split(strings.Replace(pkg, "-", "_", -1), ".")
		for i := len(chunks) - 1; i >= 0; i-- {
			protoPackageChunks = append(protoPackageChunks, chunks[i])
		}
	}
	// com.github.foo.bar.baz.qux
	protoPackage := strings.Join(protoPackageChunks, ".")

	return map[string]interface{}{
		"importPath":      importPath,
		"path":            path,
		"name":            name,
		"serviceName":     serviceName,
		"packagePath":     packagePath,
		"packageName":     packageName,
		"pbgoPackagePath": pbgoPackagePath,
		"pbgoPackageName": pbgoPackageName,
		"protoPackage":    protoPackage,
	}, nil
}
