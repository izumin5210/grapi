package generator

import (
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

// New creates a module.Generator instance.
func New(fs afero.Fs, ui cli.UI, rootDir, protoDir, protoOutDir, serverDir, pkgName string) module.Generator {
	return &generator{
		ProjectGenerator: newProjectGenerator(fs, ui),
		ServiceGenerator: newServiceGenerator(fs, ui, rootDir, protoDir, protoOutDir, serverDir, pkgName),
		CommandGenerator: newCommandGenerator(fs, ui, rootDir),
	}
}

type generator struct {
	module.ProjectGenerator
	module.ServiceGenerator
	module.CommandGenerator
}
