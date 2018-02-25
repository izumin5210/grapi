package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newProtocCommand(cfg grapicmd.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "protoc",
		Short: "Run protoc",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: force rebuild plugins option
			binDir := filepath.Join(cfg.RootDir(), "bin")
			if ok, err := afero.DirExists(cfg.Fs(), binDir); err != nil {
				return errors.Wrapf(err, "failed to get %q directory", binDir)
			} else if !ok {
				if err = cfg.Fs().MkdirAll(binDir, 0755); err != nil {
					return errors.Wrapf(err, "failed to create %q directory", binDir)
				}
			}
			for _, plugin := range cfg.ProtocConfig().Plugins {
				binName := filepath.Base(plugin.Path)
				binPath := filepath.Join(binDir, binName)
				if ok, err := afero.Exists(cfg.Fs(), binPath); err != nil {
					return errors.Wrapf(err, "failed to get %q binary", binName)
				} else if ok {
					continue
				}
				c := exec.Command("go", "install", ".")
				c.Dir = filepath.Join(cfg.RootDir(), plugin.Path)
				c.Env = append(c.Env, os.Environ()...)
				c.Env = append(c.Env, "GOBIN="+binDir)
				out, err := c.CombinedOutput()
				fmt.Println(string(out)) // FIXME
				if err != nil {
					return errors.Wrapf(err, "failed to execute command: %#v", c)
				}
			}
			return afero.Walk(
				cfg.Fs(),
				filepath.Join(cfg.RootDir(), cfg.ProtocConfig().ProtosDir),
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return errors.WithStack(err)
					}
					if !info.IsDir() && filepath.Ext(path) == ".proto" {
						outDir, err := cfg.ProtocConfig().OutDirOf(cfg.RootDir(), path)
						if err != nil {
							return errors.WithStack(err)
						}
						if ok, err := afero.DirExists(cfg.Fs(), outDir); err != nil {
							return errors.WithStack(err)
						} else if !ok {
							if err = cfg.Fs().MkdirAll(outDir, 0755); err != nil {
								return errors.Wrapf(err, "failed to create %q directory", outDir)
							}
						}
						cmds, err := cfg.ProtocConfig().Commands(cfg.RootDir(), path)
						if err != nil {
							return errors.WithStack(err)
						}
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
					}
					return nil
				},
			)
		},
	}
}
