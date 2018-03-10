package script

import (
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

// NewFactory creates a new ScriptFactory instance.
func NewFactory(fs afero.Fs, commandFactory module.CommandFactory, rootDir string) module.ScriptFactory {
	return &scriptFactory{
		fs:             fs,
		commandFactory: commandFactory,
		rootDir:        rootDir,
	}
}

type scriptFactory struct {
	fs             afero.Fs
	commandFactory module.CommandFactory
	rootDir        string
}

func (f *scriptFactory) Create(entryFilePath string) module.Script {
	return &script{
		fs:             f.fs,
		commandFactory: f.commandFactory,
		rootDir:        f.rootDir,
		entryFilePath:  entryFilePath,
	}
}
