package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
)

func newVersionCommand(ctx *grapicmd.Ctx) *cobra.Command {
	return &cobra.Command{
		Use:           "version",
		Short:         "Print version information",
		Long:          "Print version information",
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, _ []string) {
			buf := bytes.NewBufferString(ctx.AppName + " " + ctx.Version)
			if ctx.Prebuilt {
				buf.WriteString(" (" + ctx.BuildDate + " " + ctx.Revision + ")")
			}
			buf.WriteString("\n")
			fmt.Fprintf(ctx.OutWriter, buf.String())
		},
	}
}
