package usecase

import (
	"context"
	"os"
	"path/filepath"

	"github.com/izumin5210/clig/pkg/clib"
	"github.com/izumin5210/execx"
	"github.com/izumin5210/gex"
	"github.com/izumin5210/gex/pkg/tool"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/gencmd"
	_ "github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase/template"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// InitializeProjectUsecase is an interface to create a new grapi project.
type InitializeProjectUsecase interface {
	Perform(rootDir string, cfg InitConfig) error
	GenerateProject(rootDir, pkgName string) error
	InstallDeps(rootDir string, cfg InitConfig) error
}

// NewInitializeProjectUsecase creates a new InitializeProjectUsecase instance.
func NewInitializeProjectUsecase(ui cli.UI, fs afero.Fs, generator gencmd.Generator, io *clib.IO, exec *execx.Executor, gexCfg *gex.Config) InitializeProjectUsecase {
	return &initializeProjectUsecase{
		ui:        ui,
		fs:        fs,
		generator: generator,
		io:        io,
		exec:      exec,
		gexCfg:    gexCfg,
	}
}

type initializeProjectUsecase struct {
	ui        cli.UI
	fs        afero.Fs
	generator gencmd.Generator
	io        *clib.IO
	exec      *execx.Executor
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
	importPath, err := fs.GetImportPath(rootDir)
	if err != nil {
		return errors.WithStack(err)
	}

	if pkgName == "" {
		pkgName, err = fs.GetPackageName(rootDir)
		if err != nil {
			return errors.Wrap(err, "failed to decide a package name")
		}
	}

	data := map[string]interface{}{
		"packageName": pkgName,
		"importPath":  importPath,
	}
	return errors.WithStack(u.generator.Generate(data))
}

func (u *initializeProjectUsecase) InstallDeps(rootDir string, cfg InitConfig) error {
	invoke := func(name string, args ...string) error {
		cmd := u.exec.CommandContext(context.TODO(), name, args...)
		cmd.Stdout = u.io.Out
		cmd.Stderr = u.io.Err
		cmd.Stdin = u.io.In
		cmd.Dir = rootDir
		if !cfg.Dep {
			cmd.Env = append(os.Environ(), "GO111MODULE=on")
		}

		return errors.WithStack(cmd.Run())
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
		err := invoke("go", "mod", "init", filepath.Base(rootDir))
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if cfg.Dep {
		if spec := cfg.BuildSpec(); spec != "" {
			u.ui.ItemFailure("--version, --revision, --branch and --HEAD are not supported in dep mode")
		}
	} else {
		if cfg.GrapiReplacementURL != "" {
			err := invoke("go", "mod", "edit", "-replace", "github.com/izumin5210/grapi=" + cfg.GrapiReplacementURL)
			if err != nil {
				return errors.WithStack(err)
			}
		} else if spec := cfg.BuildSpec(); spec != "" {
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

	{
		// regen manifest
		path := filepath.Join(u.gexCfg.RootDir, u.gexCfg.ManifestName)
		m, err := tool.NewParser(u.fs, u.gexCfg.ManagerType).Parse(path)
		if err != nil {
			return errors.Wrapf(err, "%s was not found", path)
		}
		err = tool.NewWriter(u.fs).Write(path, m)
		if err != nil {
			return errors.WithStack(err)
		}
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
