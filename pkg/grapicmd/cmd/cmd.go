package cmd

import (
	"github.com/izumin5210/clicontrib/pkg/ccmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
)

// NewGrapiCommand creates a new command object.
func NewGrapiCommand(ctx *grapicmd.Ctx) *cobra.Command {
	var cfgFile string

	cmd := &cobra.Command{
		Use:           ctx.AppName,
		Short:         "JSON API framework implemented with gRPC and Gateway",
		Long:          "",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return errors.WithStack(ctx.Load(cfgFile))
		},
	}

	ctx.Init()

	ccmd.HandleLogFlags(cmd)

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "./"+ctx.AppName+".toml", "config file")

	cmd.AddCommand(
		newInitCommand(ctx),
		newGenerateCommand(ctx),
		newDestroyCommand(ctx),
		newProtocCommand(ctx),
		newBuildCommand(ctx),
		newVersionCommand(ctx),
	)
	cmd.AddCommand(newUserDefinedCommands(ctx)...)

	return cmd
}
