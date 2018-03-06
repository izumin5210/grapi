package usecase

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/serenize/snaker"

	"github.com/izumin5210/grapi/pkg/grapicmd/generate"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate/template"
	"github.com/izumin5210/grapi/pkg/grapicmd/ui"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// GenerateServiceUsecase is an useecase interface for geenrating .proto file and its implementation skeleton.
type GenerateServiceUsecase interface {
	Perform(path string) error
}

type generateServiceUsecase struct {
	ui        ui.UI
	generator generate.Generator
	rootDir   string
}

// NewGenerateServiceUsecase returns an new GenerateServiceUsecase implementation instance.
func NewGenerateServiceUsecase(ui ui.UI, generator generate.Generator, rootDir string) GenerateServiceUsecase {
	return &generateServiceUsecase{
		ui:        ui,
		generator: generator,
		rootDir:   rootDir,
	}
}

func (u *generateServiceUsecase) Perform(path string) error {
	// github.com/foo/bar
	importPath, err := fs.GetImportPath(u.rootDir)
	if err != nil {
		return errors.WithStack(err)
	}

	// path => baz/qux/quux

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
		pbgoPackagePath = "api_pb"
		pbgoPackageName = pbgoPackagePath
	}

	protoPackageChunks := []string{}
	for _, pkg := range strings.Split(filepath.Join(importPath, "api", filepath.Dir(path)), "/") {
		chunks := strings.Split(pkg, ".")
		for i := len(chunks) - 1; i >= 0; i-- {
			protoPackageChunks = append(protoPackageChunks, chunks[i])
		}
	}
	// com.github.foo.bar.baz.qux
	protoPackage := strings.Join(protoPackageChunks, ".")

	u.ui.Section("Generate service")

	data := map[string]interface{}{
		"importPath":      importPath,
		"path":            path,
		"name":            name,
		"serviceName":     serviceName,
		"packagePath":     packagePath,
		"packageName":     packageName,
		"pbgoPackagePath": pbgoPackagePath,
		"pbgoPackageName": pbgoPackageName,
		"protoPackage":    protoPackage,
	}
	return u.generator.Run(template.Service, data)
}
