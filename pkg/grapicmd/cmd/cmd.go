package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/clig/pkg/clib"
	"github.com/izumin5210/grapi/pkg/grapicmd"
)

// NewGrapiCommand creates a new command object.
func NewGrapiCommand(ctx *grapicmd.Ctx) *cobra.Command {
	initErr := ctx.Init()

	cmd := &cobra.Command{
		Use:           ctx.Build.AppName,
		Short:         "JSON API framework implemented with gRPC and Gateway",
		Long:          "",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return errors.WithStack(initErr)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			clib.Close()
		},
	}

	clib.AddLoggingFlags(cmd)
	clib.SetIO(cmd, ctx.IO)

	cmd.AddCommand(
		newInitCommand(ctx),
		newProtocCommand(ctx),
		newBuildCommand(ctx),
		clib.NewVersionCommand(ctx.Build),
	)
	cmd.AddCommand(newGenerateCommands(ctx)...)
	cmd.AddCommand(newUserDefinedCommands(ctx)...)

	return cmd
}
