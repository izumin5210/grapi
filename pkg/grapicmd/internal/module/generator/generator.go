package generator

import (
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

// New creates a module.Generator instance.
func New(fs afero.Fs, ui module.UI, rootDir, protoDir, protoOutDir, serverDir, version string) module.Generator {
	return &generator{
		ProjectGenerator: newProjectGenerator(fs, ui, version),
		ServiceGenerator: newServiceGenerator(fs, ui, rootDir, protoDir, protoOutDir, serverDir),
		CommandGenerator: newCommandGenerator(fs, ui, rootDir),
	}
}

type generator struct {
	module.ProjectGenerator
	module.ServiceGenerator
	module.CommandGenerator
}
