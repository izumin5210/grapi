package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/command"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

func newUserDefinedCommand(cfg grapicmd.Config, rootDir, entryPath string) *cobra.Command {
	name := filepath.Base(filepath.Dir(entryPath))
	binDir := filepath.Join(rootDir, "bin")
	binPath := filepath.Join(binDir, name)
	return &cobra.Command{
		Use:           name,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(c *cobra.Command, args []string) error {
			executor := command.NewExecutor(rootDir, cfg.OutWriter(), cfg.ErrWriter(), cfg.InReader())

			err := fs.CreateDirIfNotExists(cfg.Fs(), binDir)
			if err != nil {
				return errors.WithStack(err)
			}

			out, err := executor.Exec([]string{"go", "build", "-v", "-o=" + binPath, entryPath})
			if err != nil {
				fmt.Println(string(out))
				return errors.Wrapf(err, "failed to build %q", entryPath)
			}

			_, err = executor.Exec([]string{binPath}, command.WithIOConnected())
			return errors.WithStack(err)
		},
	}
}
