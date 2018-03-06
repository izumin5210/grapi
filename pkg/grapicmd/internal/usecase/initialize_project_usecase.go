package usecase

import (
	"github.com/pkg/errors"

	"github.com/izumin5210/grapi/pkg/grapicmd/command"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate/template"
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
func NewInitializeProjectUsecase(ui module.UI, generator generate.Generator, executor command.Executor) InitializeProjectUsecase {
	return &initializeProjectUsecase{
		ui:        ui,
		generator: generator,
		executor:  executor,
	}
}

type initializeProjectUsecase struct {
	ui        module.UI
	generator generate.Generator
	executor  command.Executor
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
	return errors.WithStack(u.generator.Run(template.Init, data))
}

func (u *initializeProjectUsecase) InstallDeps(rootDir string) error {
	_, err := u.executor.Exec(
		[]string{"dep", "ensure", "-v"},
		command.WithIOConnected(),
		command.WithDir(rootDir),
	)
	return errors.WithStack(err)
}
