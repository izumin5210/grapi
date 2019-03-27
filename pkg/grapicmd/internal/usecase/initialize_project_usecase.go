package usecase

import (
	"context"

	"github.com/izumin5210/gex"
	"github.com/pkg/errors"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/excmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

// InitializeProjectUsecase is an interface to create a new grapi project.
type InitializeProjectUsecase interface {
	Perform(rootDir string, cfg InitConfig) error
	GenerateProject(rootDir, pkgName string) error
	InstallDeps(rootDir string, cfg InitConfig) error
}

// NewInitializeProjectUsecase creates a new InitializeProjectUsecase instance.
func NewInitializeProjectUsecase(ui cli.UI, generator module.ProjectGenerator, excmd excmd.Executor, gexCfg *gex.Config) InitializeProjectUsecase {
	return &initializeProjectUsecase{
		ui:        ui,
		generator: generator,
		excmd:     excmd,
		gexCfg:    gexCfg,
	}
}

type initializeProjectUsecase struct {
	ui        cli.UI
	generator module.ProjectGenerator
	excmd     excmd.Executor
	gexCfg    *gex.Config
}

func (u *initializeProjectUsecase) Perform(rootDir string, cfg InitConfig) error {
	u.ui.Section("Initialize project")

	var err error
	err = u.GenerateProject(rootDir, cfg.Package)
	if err != nil {
		return errors.Wrap(err, "failed to initialize project")
	}

	u.ui.Subsection("Install dependencies")
	err = u.InstallDeps(rootDir, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to execute `dep ensure`")
	}

	return nil
}

func (u *initializeProjectUsecase) GenerateProject(rootDir, pkgName string) error {
	return errors.WithStack(u.generator.GenerateProject(rootDir, pkgName))
}

func (u *initializeProjectUsecase) InstallDeps(rootDir string, cfg InitConfig) error {
	var (
		name string
		args []string
		opts = []excmd.Option{
			excmd.WithDir(rootDir),
			excmd.WithIOConnected(),
		}
	)
	if cfg.Dep {
		name = "dep"
		args = []string{"init"}
	} else {
		name = "go"
		args = []string{"mod", "init"}
		opts = append(opts, excmd.WithEnv("GO111MODULE", "on"))
	}
	_, err := u.excmd.Exec(context.Background(), name, append([]excmd.Option{excmd.WithArgs(args...)}, opts...)...)
	if err != nil {
		return errors.WithStack(err)
	}

	if cfg.Dep {
		if spec := cfg.BuildSpec(); spec != "" {
			u.ui.ItemFailure("--version, --revision, --branch and --HEAD are not supported in dep mode")
		}
	} else {
		if spec := cfg.BuildSpec(); spec != "" {
			pkg := "github.com/izumin5210/grapi/pkg/grapiserver"
			args = []string{"get", pkg + spec}
			_, err := u.excmd.Exec(context.Background(), name, append([]excmd.Option{excmd.WithArgs(args...)}, opts...)...)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		_, err := u.excmd.Exec(context.Background(), "go", append([]excmd.Option{excmd.WithArgs("get", "./...")}, opts...)...)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
