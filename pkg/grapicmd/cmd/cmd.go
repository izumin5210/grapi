package cmd

import (
	"github.com/izumin5210/clicontrib/pkg/ccmd"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
)

func provideGrapiCommand(
	cfg *grapicmd.Config,
	initCmd initCmd,
	generateCmd generateCmd,
	destroyCmd destroyCmd,
	protocCmd protocCmd,
	buildCmd buildCmd,
	versionCmd versionCmd,
	userDefinedCmds userDefinedCmds,
) *cobra.Command {
	var err error

	cmd := &cobra.Command{
		Use:           cfg.AppName,
		Short:         "JSON API framework implemented with gRPC and Gateway",
		Long:          "",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return err
		},
	}

	var cfgFile string
	cobra.OnInitialize(func() { cfg.Init(cfgFile) })
	ccmd.HandleLogFlags(cmd)

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "./"+cfg.AppName+".toml", "config file")

	cmd.AddCommand(
		initCmd,
		generateCmd,
		destroyCmd,
		protocCmd,
		buildCmd,
		versionCmd,
	)
	cmd.AddCommand(userDefinedCmds...)

	return cmd
}
