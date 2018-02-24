package cmd

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// NewGrapiCommand creates a new command object.
func NewGrapiCommand(fs afero.Fs, inReader io.Reader, outWriter, errWriter io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grapi",
		Short: "grapi is JSON API implemented with gRPC and Gateway",
		Long:  "",
	}
	cmd.AddCommand(newInitCommand())

	udCmds := make([]*cobra.Command, 0)
	wd, err := os.Getwd()
	if err == nil {
		paths, err := afero.Glob(fs, filepath.Join(wd, "cmd/*/run.go"))
		if err == nil {
			for _, path := range paths {
				udCmds = append(udCmds, newUserDefinedCommand(path))
			}
		}
	}
	if len(udCmds) > 0 {
		cmd.AddCommand(udCmds...)
	}

	return cmd
}
