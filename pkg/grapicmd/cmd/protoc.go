package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/izumin5210/clicontrib/clog"
	"github.com/izumin5210/grapi/pkg/grapicmd"
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
			rootDir, ok := fs.LookupRoot(cfg.Fs(), cfg.CurrentDir())
			if !ok {
				return errors.New("protoc command should be execute inside a grapi applicaiton directory")
			}

			// TODO: force rebuild plugins option
			binDir := filepath.Join(rootDir, "bin")
			if err := fs.CreateDirIfNotExists(cfg.Fs(), binDir); err != nil {
				return errors.WithStack(err)
			}
			ui.Section("Install plugins")
			var errs []error
			for _, plugin := range cfg.ProtocConfig().Plugins {
				ok, err := installPlugin(cfg.Fs(), plugin, rootDir, binDir)
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
			protoFiles, err := cfg.ProtocConfig().ProtoFiles(cfg.Fs(), rootDir)
			if err != nil {
				return errors.WithStack(err)
			}
			for _, path := range protoFiles {
				// err = executeProtoc(cfg.Fs(), cfg.ProtocConfig(), binDir, rootDir, path)
				err = executeProtoc(cfg.Fs(), cfg.ProtocConfig(), rootDir, binDir, path)
				relPath, _ := filepath.Rel(rootDir, path)
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

func installPlugin(fs afero.Fs, plugin *protoc.Plugin, rootDir, binDir string) (bool, error) {
	binPath := filepath.Join(binDir, plugin.BinName())
	if ok, err := afero.Exists(fs, binPath); err != nil {
		return false, errors.Wrapf(err, "failed to get %q binary", plugin.BinName())
	} else if ok {
		return false, nil
	}
	c := exec.Command("go", "install", ".")
	c.Dir = filepath.Join(rootDir, plugin.Path)
	if ok, _ := afero.DirExists(fs, c.Dir); !ok {
		return false, errors.Errorf("%s is not found", plugin.Path)
	}
	c.Env = append(c.Env, os.Environ()...)
	c.Env = append(c.Env, "GOBIN="+binDir)
	out, err := c.CombinedOutput()
	clog.Debug("execute", "command", c.Args, "out", string(out), "dir", c.Dir)
	if err != nil {
		return false, errors.Wrapf(err, "failed to execute command: %v, output: %s", c.Args, string(out))
	}
	return true, nil
}

func executeProtoc(afs afero.Fs, pcfg *protoc.Config, rootDir, binDir, protoPath string) error {
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
		c := exec.Command(cmd[0], cmd[1:]...)
		c.Dir = rootDir
		c.Env = append(c.Env, os.Environ()...)
		c.Env = append(c.Env, "PATH="+binDir+":"+os.Getenv("PATH"))
		out, err := c.CombinedOutput()
		clog.Debug("execute", "command", c.Args, "out", string(out), "dir", c.Dir)
		if err != nil {
			return errors.Wrapf(err, "failed to execute command: %v, output: %s", c.Args, string(out))
		}
	}
	return nil
}
