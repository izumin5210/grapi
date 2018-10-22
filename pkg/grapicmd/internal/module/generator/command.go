package generator

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator/template"
)

type commandGenerator struct {
	baseGenerator
	rootDir string
}

func newCommandGenerator(fs afero.Fs, ui cli.UI, rootDir string) module.CommandGenerator {
	return &commandGenerator{
		baseGenerator: newBaseGenerator(template.Command, fs, ui),
		rootDir:       rootDir,
	}
}

func (g *commandGenerator) GenerateCommand(name string) error {
	return errors.WithStack(g.Generate(g.rootDir, map[string]string{"name": name}, generationConfig{}))
}

func (g *commandGenerator) DestroyCommand(name string) error {
	return errors.WithStack(g.Destroy(g.rootDir, map[string]string{"name": name}))
}
