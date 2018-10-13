package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
)

func newVersionCommand(cfg *grapicmd.Config) *cobra.Command {
	return &cobra.Command{
		Use:           "version",
		Short:         "Print version information",
		Long:          "Print version information",
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, _ []string) {
			buf := bytes.NewBufferString(cfg.AppName + " " + cfg.Version)
			if cfg.Prebuilt {
				buf.WriteString(" (" + cfg.BuildDate + " " + cfg.Revision + ")")
			}
			buf.WriteString("\n")
			fmt.Fprintf(cfg.OutWriter, buf.String())
		},
	}
}
