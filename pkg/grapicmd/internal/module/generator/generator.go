package generator

import (
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

// New creates a module.Generator instance.
func New(fs afero.Fs, ui cli.UI) module.Generator {
	return &generator{
		ProjectGenerator: newProjectGenerator(fs, ui),
	}
}

type generator struct {
	module.ProjectGenerator
}
