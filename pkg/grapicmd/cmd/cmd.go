package cmd

import (
	"github.com/izumin5210/clicontrib/pkg/ccmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
)

// NewGrapiCommand creates a new command object.
func NewGrapiCommand(cfg *grapicmd.Config) *cobra.Command {
	var cfgFile string

	cmd := &cobra.Command{
		Use:           cfg.AppName,
		Short:         "JSON API framework implemented with gRPC and Gateway",
		Long:          "",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return errors.WithStack(cfg.Init(cfgFile))
		},
	}

	ccmd.HandleLogFlags(cmd)

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "./"+cfg.AppName+".toml", "config file")

	cmd.AddCommand(
		newInitCommand(cfg),
		newGenerateCommand(cfg),
		newDestroyCommand(cfg),
		newProtocCommand(cfg),
		newBuildCommand(cfg),
		newVersionCommand(cfg),
	)
	cmd.AddCommand(newUserDefinedCommands(cfg)...)

	return cmd
}
