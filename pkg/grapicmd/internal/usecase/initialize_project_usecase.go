package usecase

import (
	"github.com/pkg/errors"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

// InitializeProjectUsecase is an interface to create a new grapi project.
type InitializeProjectUsecase interface {
	Perform(rootDir string, depSkipped, headUsed bool) error
	GenerateProject(rootDir string, headUsed bool) error
	InstallDeps(rootDir string) error
}

// NewInitializeProjectUsecase creates a new InitializeProjectUsecase instance.
func NewInitializeProjectUsecase(ui module.UI, generator module.ProjectGenerator, commandFactory module.CommandFactory, version string) InitializeProjectUsecase {
	return &initializeProjectUsecase{
		ui:             ui,
		generator:      generator,
		commandFactory: commandFactory,
		version:        version,
	}
}

type initializeProjectUsecase struct {
	ui             module.UI
	generator      module.ProjectGenerator
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
	return errors.WithStack(u.generator.GenerateProject(rootDir, headUsed))
}

func (u *initializeProjectUsecase) InstallDeps(rootDir string) error {
	cmd := u.commandFactory.Create([]string{"dep", "ensure", "-v"})
	_, err := cmd.ConnectIO().SetDir(rootDir).Exec()
	return errors.WithStack(err)
}
