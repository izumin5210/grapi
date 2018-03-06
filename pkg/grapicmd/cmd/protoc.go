package cmd

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/command"
	"github.com/izumin5210/grapi/pkg/grapicmd/protoc"
	"github.com/izumin5210/grapi/pkg/grapicmd/ui"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

func newProtocCommand(cfg grapicmd.Config, ui ui.UI) *cobra.Command {
	return &cobra.Command{
		Use:           "protoc",
		Short:         "Run protoc",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.IsInsideApp() {
				return errors.New("protoc command should be execute inside a grapi applicaiton directory")
			}
			executor := command.NewExecutor(cfg.RootDir(), cfg.OutWriter(), cfg.ErrWriter(), cfg.InReader())

			// TODO: force rebuild plugins option
			binDir := filepath.Join(cfg.RootDir(), "bin")
			if err := fs.CreateDirIfNotExists(cfg.Fs(), binDir); err != nil {
				return errors.WithStack(err)
			}
			ui.Section("Install plugins")
			var errs []error
			for _, plugin := range cfg.ProtocConfig().Plugins {
				ok, err := installPlugin(cfg.Fs(), executor, plugin, cfg.RootDir(), binDir)
				if err != nil {
					errs = append(errs, err)
					ui.ItemFailure(plugin.BinName())
				} else if !ok {
					ui.ItemSkipped(plugin.BinName())
				} else {
					ui.ItemSuccess(plugin.BinName())
				}
			}
			if len(errs) > 0 {
				for _, err := range errs {
					ui.Error(err.Error())
				}
				return errors.New("failed to install protoc plugins")
			}

			ui.Section("Execute protoc")
			protoFiles, err := cfg.ProtocConfig().ProtoFiles(cfg.Fs(), cfg.RootDir())
			if err != nil {
				return errors.WithStack(err)
			}
			for _, path := range protoFiles {
				err = executeProtoc(cfg.Fs(), executor, cfg.ProtocConfig(), cfg.RootDir(), binDir, path)
				relPath, _ := filepath.Rel(cfg.RootDir(), path)
				if err == nil {
					ui.ItemSuccess(relPath)
				} else {
					errs = append(errs, err)
					ui.ItemFailure(relPath)
				}
			}
			if len(errs) > 0 {
				for _, err := range errs {
					ui.Error(err.Error())
				}
				return errors.New("failed to execute protoc")
			}
			return nil
		},
	}
}

func installPlugin(fs afero.Fs, executor command.Executor, plugin *protoc.Plugin, rootDir, binDir string) (bool, error) {
	binPath := filepath.Join(binDir, plugin.BinName())
	if ok, err := afero.Exists(fs, binPath); err != nil {
		return false, errors.Wrapf(err, "failed to get %q binary", plugin.BinName())
	} else if ok {
		return false, nil
	}
	dir := filepath.Join(rootDir, plugin.Path)
	if ok, _ := afero.DirExists(fs, dir); !ok {
		return false, errors.Errorf("%s is not found", plugin.Path)
	}
	out, err := executor.Exec(
		[]string{"go", "install", "."},
		command.WithDir(dir),
		command.WithEnv("GOBIN", binDir),
	)
	if err != nil {
		return false, errors.Wrapf(err, "failed to execute command: %s", string(out))
	}
	return true, nil
}

func executeProtoc(afs afero.Fs, executor command.Executor, pcfg *protoc.Config, rootDir, binDir, protoPath string) error {
	outDir, err := pcfg.OutDirOf(rootDir, protoPath)
	if err != nil {
		return errors.WithStack(err)
	}
	if err = fs.CreateDirIfNotExists(afs, outDir); err != nil {
		return errors.WithStack(err)
	}
	cmds, err := pcfg.Commands(rootDir, protoPath)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, cmd := range cmds {
		out, err := executor.Exec(cmd, command.WithEnv("PATH", binDir+":"+os.Getenv("PATH")))
		if err != nil {
			return errors.Wrapf(err, "failed to execute command: %s", string(out))
		}
	}
	return nil
}
