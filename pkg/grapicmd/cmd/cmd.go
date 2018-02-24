package cmd

import (
	"io"

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
	return cmd
}
