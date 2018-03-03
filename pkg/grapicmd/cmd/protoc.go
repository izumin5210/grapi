package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
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
				return errors.New("protoc command should execut inside a grapi applicaiton directory")
			}

			// TODO: force rebuild plugins option
			binDir := filepath.Join(rootDir, "bin")
			if err := fs.CreateDirIfNotExists(cfg.Fs(), binDir); err != nil {
				return errors.WithStack(err)
			}
			ui.Section("Install plugins")
			for _, plugin := range cfg.ProtocConfig().Plugins {
				binName := filepath.Base(plugin.Path)
				binPath := filepath.Join(binDir, binName)
				if ok, err := afero.Exists(cfg.Fs(), binPath); err != nil {
					return errors.Wrapf(err, "failed to get %q binary", binName)
				} else if ok {
					ui.ItemSkipped(filepath.Base(plugin.Path))
					continue
				}
				c := exec.Command("go", "install", ".")
				c.Dir = filepath.Join(rootDir, plugin.Path)
				c.Env = append(c.Env, os.Environ()...)
				c.Env = append(c.Env, "GOBIN="+binDir)
				out, err := c.CombinedOutput()
				if err != nil {
					fmt.Println(string(out)) // FIXME
					return errors.Wrapf(err, "failed to execute command: %#v", c)
				}
				ui.ItemSuccess(filepath.Base(plugin.Path))
			}
			ui.Section("Execute protoc")
			return afero.Walk(
				cfg.Fs(),
				filepath.Join(rootDir, cfg.ProtocConfig().ProtosDir),
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return errors.WithStack(err)
					}
					if !info.IsDir() && filepath.Ext(path) == ".proto" {
						outDir, err := cfg.ProtocConfig().OutDirOf(rootDir, path)
						if err != nil {
							return errors.WithStack(err)
						}
						if err = fs.CreateDirIfNotExists(cfg.Fs(), outDir); err != nil {
							return errors.WithStack(err)
						}
						cmds, err := cfg.ProtocConfig().Commands(rootDir, path)
						if err != nil {
							return errors.WithStack(err)
						}
						relPath, _ := filepath.Rel(rootDir, path)
						for _, cmd := range cmds {
							c := exec.Command(cmd[0], cmd[1:]...)
							c.Env = append(c.Env, os.Environ()...)
							c.Env = append(c.Env, "PATH="+binDir+":"+os.Getenv("PATH"))
							out, err := c.CombinedOutput()
							if err != nil {
								fmt.Println(string(out)) // FIXME
								return errors.Wrapf(err, "failed to execute command: %#v", c)
							}
						}
						ui.ItemSuccess(relPath)
					}
					return nil
				},
			)
		},
	}
}
