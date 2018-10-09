package usecase

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/izumin5210/gex"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/protoc"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// ExecuteProtocUsecase is an useecase interface for executing protoc module.
type ExecuteProtocUsecase interface {
	Perform() error
	InstallPlugins() error
	ExecuteProtoc() error
}

type executeProtocUsecase struct {
	cfg            *protoc.Config
	fs             afero.Fs
	ui             module.UI
	commandFactory module.CommandFactory
	gexCfg         *gex.Config
	rootDir        string
}

// NewExecuteProtocUsecase returns an new ExecuteProtocUsecase implementation instance.
func NewExecuteProtocUsecase(cfg *protoc.Config, fs afero.Fs, ui module.UI, commandFactory module.CommandFactory, gexCfg *gex.Config, rootDir string) ExecuteProtocUsecase {
	return &executeProtocUsecase{
		cfg:            cfg,
		fs:             fs,
		ui:             ui,
		commandFactory: commandFactory,
		gexCfg:         gexCfg,
		rootDir:        rootDir,
	}
}

func (u *executeProtocUsecase) Perform() error {
	u.ui.Section("Execute protoc")
	u.ui.Subsection("Install plugins")
	err := errors.WithStack(u.InstallPlugins())
	if err != nil {
		return err
	}
	u.ui.Subsection("Execute protoc")
	return errors.WithStack(u.ExecuteProtoc())
}

func (u *executeProtocUsecase) InstallPlugins() error {
	repo, err := u.gexCfg.Create()

	if err == nil {
		err = repo.BuildAll(context.TODO())
	}

	return errors.WithStack(err)
}

func (u *executeProtocUsecase) ExecuteProtoc() error {
	protoFiles, err := u.cfg.ProtoFiles(u.fs, u.rootDir)
	if err != nil {
		return errors.WithStack(err)
	}
	var errs []error
	for _, path := range protoFiles {
		err = u.executeProtoc(path)
		relPath, _ := filepath.Rel(u.rootDir, path)
		if err == nil {
			u.ui.ItemSuccess(relPath)
		} else {
			errs = append(errs, err)
			u.ui.ItemFailure(relPath)
		}
	}
	if len(errs) > 0 {
		for _, err := range errs {
			u.ui.Error(err.Error())
		}
		return errors.New("failed to execute protoc")
	}
	return nil
}

func (u *executeProtocUsecase) executeProtoc(protoPath string) error {
	outDir, err := u.cfg.OutDirOf(u.rootDir, protoPath)
	if err != nil {
		return errors.WithStack(err)
	}
	if err = fs.CreateDirIfNotExists(u.fs, outDir); err != nil {
		return errors.WithStack(err)
	}
	cmds, err := u.cfg.Commands(u.rootDir, protoPath)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, cmd := range cmds {
		path := strings.Join([]string{
			filepath.Join(u.rootDir, u.gexCfg.BinDirName),
			os.Getenv("PATH"),
		}, string(filepath.ListSeparator))
		out, err := u.commandFactory.Create(cmd).AddEnv("PATH", path).SetDir(u.rootDir).Exec()
		if err != nil {
			return errors.Wrapf(err, "failed to execute module: %s", string(out))
		}
	}
	return nil
}
