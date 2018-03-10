package usecase

import (
	"github.com/pkg/errors"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// InitializeProjectUsecase is an interface to create a new grapi project.
type InitializeProjectUsecase interface {
	Perform(rootDir string, depSkipped bool) error
	GenerateProject(rootDir string) error
	InstallDeps(rootDir string) error
}

// NewInitializeProjectUsecase creates a new InitializeProjectUsecase instance.
func NewInitializeProjectUsecase(ui module.UI, generator module.Generator, commandFactory module.CommandFactory) InitializeProjectUsecase {
	return &initializeProjectUsecase{
		ui:             ui,
		generator:      generator,
		commandFactory: commandFactory,
	}
}

type initializeProjectUsecase struct {
	ui             module.UI
	generator      module.Generator
	commandFactory module.CommandFactory
}

func (u *initializeProjectUsecase) Perform(rootDir string, depSkipped bool) error {
	u.ui.Section("Initialize project")

	var err error
	err = u.GenerateProject(rootDir)
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

func (u *initializeProjectUsecase) GenerateProject(rootDir string) error {
	importPath, err := fs.GetImportPath(rootDir)
	if err != nil {
		return errors.WithStack(err)
	}
	data := map[string]string{
		"importPath": importPath,
	}
	return errors.WithStack(u.generator.Exec(rootDir, data))
}

func (u *initializeProjectUsecase) InstallDeps(rootDir string) error {
	cmd := u.commandFactory.Create([]string{"dep", "ensure", "-v"})
	_, err := cmd.ConnectIO().SetDir(rootDir).Exec()
	return errors.WithStack(err)
}
