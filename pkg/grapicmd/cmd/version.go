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
			b := ctx.Build
			buf := bytes.NewBufferString(b.AppName + " " + b.Version)
			if b.Prebuilt {
				buf.WriteString(" (" + b.BuildDate + " " + b.Revision + ")")
			}
			buf.WriteString("\n")
			fmt.Fprintf(ctx.IO.Out, buf.String())
		},
	}
}
