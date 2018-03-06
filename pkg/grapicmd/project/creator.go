package project

import (
	"github.com/pkg/errors"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/command"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate/template"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// Creator is an interface to create a new grapi project.
type Creator interface {
	Run() error
}

// Config contains configurations for creating a project.
type Config struct {
	grapicmd.Config
	RootDir    string
	DepSkipped bool
}

// NewCreator creates a new Creator instance.
func NewCreator(ui module.UI, generator generate.Generator, executor command.Executor, cfg *Config) Creator {
	return &creator{
		ui:        ui,
		generator: generator,
		executor:  executor,
		cfg:       cfg,
	}
}

type creator struct {
	cfg       *Config
	ui        module.UI
	generator generate.Generator
	executor  command.Executor
}

func (c *creator) Run() error {
	var err error
	err = c.initProject()
	if err != nil {
		return errors.Wrap(err, "failed to initialize project")
	}

	if !c.cfg.DepSkipped {
		err = c.installDeps()
		if err != nil {
			return errors.Wrap(err, "failed to execute `dep ensure`")
		}
	}

	return nil
}

func (c *creator) initProject() error {
	c.ui.Section("Initialize project")
	importPath, err := fs.GetImportPath(c.cfg.RootDir)
	if err != nil {
		return errors.WithStack(err)
	}
	data := map[string]string{
		"importPath": importPath,
	}
	return errors.WithStack(c.generator.Run(template.Init, data))
}

func (c *creator) installDeps() error {
	c.ui.Subsection("Install dependencies")
	_, err := c.executor.Exec(
		[]string{"dep", "ensure", "-v"},
		command.WithIOConnected(),
	)
	return errors.WithStack(err)
}
