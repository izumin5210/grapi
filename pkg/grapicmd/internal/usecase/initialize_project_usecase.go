package usecase

import (
	"github.com/pkg/errors"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// InitializeProjectUsecase is an interface to create a new grapi project.
type InitializeProjectUsecase interface {
	Perform(rootDir string, depSkipped, headUsed bool) error
	GenerateProject(rootDir string, headUsed bool) error
	InstallDeps(rootDir string) error
}

// NewInitializeProjectUsecase creates a new InitializeProjectUsecase instance.
func NewInitializeProjectUsecase(ui module.UI, generator module.Generator, commandFactory module.CommandFactory, version string) InitializeProjectUsecase {
	return &initializeProjectUsecase{
		ui:             ui,
		generator:      generator,
		commandFactory: commandFactory,
		version:        version,
	}
}

type initializeProjectUsecase struct {
	ui             module.UI
	generator      module.Generator
	commandFactory module.CommandFactory
	version        string
}

func (u *initializeProjectUsecase) Perform(rootDir string, depSkipped, headUsed bool) error {
	u.ui.Section("Initialize project")

	var err error
	err = u.GenerateProject(rootDir, headUsed)
	if err != nil {
		return errors.Wrap(err, "failed to initialize project")
	}

	u.ui.Subsection("Install dependencies")
	if !depSkipped {
		err = u.InstallDeps(rootDir)
		if err != nil {
			return errors.Wrap(err, "failed to execute `dep ensure`")
		}
	}

	return nil
}

func (u *initializeProjectUsecase) GenerateProject(rootDir string, headUsed bool) error {
	importPath, err := fs.GetImportPath(rootDir)
	if err != nil {
		return errors.WithStack(err)
	}
	data := map[string]interface{}{
		"importPath": importPath,
		"version":    u.version,
		"headUsed":   headUsed,
	}
	return errors.WithStack(u.generator.Exec(rootDir, data))
}

func (u *initializeProjectUsecase) InstallDeps(rootDir string) error {
	cmd := u.commandFactory.Create([]string{"dep", "ensure", "-v"})
	_, err := cmd.ConnectIO().SetDir(rootDir).Exec()
	return errors.WithStack(err)
}
