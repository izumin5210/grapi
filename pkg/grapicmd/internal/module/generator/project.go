package generator

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator/template"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

type projectGenerator struct {
	baseGenerator
	version string
}

func newProjectGenerator(fs afero.Fs, ui module.UI, version string) module.ProjectGenerator {
	return &projectGenerator{
		baseGenerator: newBaseGenerator(template.Init, fs, ui),
		version:       version,
	}
}

func (g *projectGenerator) GenerateProject(rootDir, pkgName string, useHead bool) error {
	importPath, err := fs.GetImportPath(rootDir)
	if err != nil {
		return errors.WithStack(err)
	}

	if pkgName == "" {
		pkgName, err = fs.GetPackageName(rootDir)
		if err != nil {
			return errors.Wrap(err, "failed to decide a package name")
		}
	}

	data := map[string]interface{}{
		"packageName": pkgName,
		"importPath":  importPath,
		"version":     g.version,
		"headUsed":    useHead,
	}
	return g.Generate(rootDir, data, generationConfig{})
}
