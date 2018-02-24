package cmd

import (
	"os/exec"
	"path/filepath"
	"plugin"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newUserDefinedCommand(entryPath string) *cobra.Command {
	dirPath := filepath.Dir(entryPath)
	name := filepath.Base(dirPath)
	soPath := filepath.Join(dirPath, name+".so")
	return &cobra.Command{
		Use: name,
		RunE: func(c *cobra.Command, args []string) error {
			err := exec.Command("go", "build", "-v", "-buildmode=plugin", "-o="+soPath, entryPath).Run()
			if err != nil {
				return errors.Wrapf(err, "failed to build %q", entryPath)
			}

			p, err := plugin.Open(soPath)
			if err != nil {
				return errors.Wrap(err, "failed to laod plugin")
			}

			runSym, err := p.Lookup("Run")
			if err != nil {
				return errors.Wrap(err, "failed to lookup `func Run(args []string) error`")
			}

			run, ok := runSym.(func([]string) error)
			if !ok {
				return errors.Wrap(err, "`Run` should be `func(args []string) error`")
			}

			return run(args)
		},
	}
}
