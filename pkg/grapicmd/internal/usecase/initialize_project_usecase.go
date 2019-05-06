package usecase

import (
	"context"
	"path/filepath"

	"github.com/izumin5210/gex"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

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
func NewInitializeProjectUsecase(ui cli.UI, fs afero.Fs, generator module.ProjectGenerator, excmd excmd.Executor, gexCfg *gex.Config) InitializeProjectUsecase {
	return &initializeProjectUsecase{
		ui:        ui,
		fs:        fs,
		generator: generator,
		excmd:     excmd,
		gexCfg:    gexCfg,
	}
}

type initializeProjectUsecase struct {
	ui        cli.UI
	fs        afero.Fs
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
		return errors.Wrap(err, "failed to install dependencies")
	}

	return nil
}

func (u *initializeProjectUsecase) GenerateProject(rootDir, pkgName string) error {
	return errors.WithStack(u.generator.GenerateProject(rootDir, pkgName))
}

func (u *initializeProjectUsecase) InstallDeps(rootDir string, cfg InitConfig) error {
	opts := []excmd.Option{
		excmd.WithDir(rootDir),
		excmd.WithIOConnected(),
	}
	if !cfg.Dep {
		opts = append(opts, excmd.WithEnv("GO111MODULE", "on"))
	}

	invoke := func(name string, args ...string) error {
		_, err := u.excmd.Exec(context.Background(), name, append([]excmd.Option{excmd.WithArgs(args...)}, opts...)...)
		return errors.WithStack(err)
	}

	if cfg.Dep {
		err := invoke("dep", "init")
		if err != nil {
			return errors.WithStack(err)
		}
	} else if ok, _ := afero.Exists(u.fs, filepath.Join(rootDir, "go.mod")); ok {
		err := invoke("go", "mod", "tidy")
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		err := invoke("go", "mod", "init")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if cfg.Dep {
		if spec := cfg.BuildSpec(); spec != "" {
			u.ui.ItemFailure("--version, --revision, --branch and --HEAD are not supported in dep mode")
		}
	} else {
		if spec := cfg.BuildSpec(); spec != "" {
			pkg := "github.com/izumin5210/grapi/pkg/grapiserver"
			err := invoke("go", "get", pkg+spec)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		err := invoke("go", "get", "./...")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	u.gexCfg.WorkingDir = rootDir
	u.gexCfg.RootDir = rootDir
	toolRepo, err := u.gexCfg.Create()
	if err != nil {
		return errors.WithStack(err)
	}
	err = toolRepo.BuildAll(context.Background())
	if err != nil {
		return errors.WithStack(err)
	}

	if !cfg.Dep {
		err := invoke("go", "mod", "tidy")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
