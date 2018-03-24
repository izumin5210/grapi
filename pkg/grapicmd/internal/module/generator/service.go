package generator

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/serenize/snaker"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator/template"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

type serviceGenerator struct {
	baseGenerator
	rootDir string
}

func newServiceGenerator(fs afero.Fs, ui module.UI, rootDir string) module.ServiceGenerator {
	return &serviceGenerator{
		baseGenerator: newBaseGenerator(template.Service, fs, ui),
		rootDir:       rootDir,
	}
}

func (g *serviceGenerator) GenerateService(name string) error {
	data, err := g.createParams(name)
	if err != nil {
		return errors.WithStack(err)
	}
	return g.Generate(g.rootDir, data)
}

func (g *serviceGenerator) DestroyService(name string) error {
	data, err := g.createParams(name)
	if err != nil {
		return errors.WithStack(err)
	}
	return g.Destroy(g.rootDir, data)
}

func (g *serviceGenerator) createParams(path string) (map[string]interface{}, error) {
	// github.com/foo/bar
	importPath, err := fs.GetImportPath(g.rootDir)
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
