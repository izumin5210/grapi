package cmd

import (
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
)

func newVersionCommand(cfg grapicmd.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print version information",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Printf("%s %s (%s)\n", cfg.AppName(), cfg.Version(), cfg.Revision())
		},
	}
}
