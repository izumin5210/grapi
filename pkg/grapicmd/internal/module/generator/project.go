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

func (g *projectGenerator) GenerateProject(rootDir string, useHead bool) error {
	importPath, err := fs.GetImportPath(rootDir)
	if err != nil {
		return errors.WithStack(err)
	}
	data := map[string]interface{}{
		"importPath": importPath,
		"version":    g.version,
		"headUsed":   useHead,
	}
	return g.Generate(rootDir, data)
}
