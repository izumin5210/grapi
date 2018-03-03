package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

func newUserDefinedCommand(cfg grapicmd.Config, rootDir, entryPath string) *cobra.Command {
	name := filepath.Base(filepath.Dir(entryPath))
	binDir := filepath.Join(rootDir, "bin")
	binPath := filepath.Join(binDir, name)
	return &cobra.Command{
		Use: name,
		RunE: func(c *cobra.Command, args []string) error {
			err := fs.CreateDirIfNotExists(cfg.Fs(), binDir)
			if err != nil {
				return errors.WithStack(err)
			}

			out, err := exec.Command("go", "build", "-v", "-o="+binPath, entryPath).CombinedOutput()
			if err != nil {
				fmt.Println(string(out))
				return errors.Wrapf(err, "failed to build %q", entryPath)
			}

			cmd := exec.Command(binPath)
			cmd.Stdin = cfg.InReader()
			cmd.Stdout = cfg.OutWriter()
			cmd.Stderr = cfg.ErrWriter()
			return errors.WithStack(cmd.Run())
		},
	}
}
