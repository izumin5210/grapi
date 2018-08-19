package cmd

import (
	"path/filepath"

	"github.com/izumin5210/clicontrib/pkg/ccmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
)

// NewGrapiCommand creates a new command object.
func NewGrapiCommand(cfg grapicmd.Config) *cobra.Command {
	var err error

	cmd := &cobra.Command{
		Use:           cfg.AppName(),
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

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "./"+cfg.AppName()+".toml", "config file")

	ac := di.NewAppComponent(cfg)

	cmd.AddCommand(newInitCommand(ac))
	cmd.AddCommand(newGenerateCommand(ac))
	cmd.AddCommand(newDestroyCommand(ac))
	cmd.AddCommand(newProtocCommand(ac))
	cmd.AddCommand(newBuildCommand(ac))
	cmd.AddCommand(newVersionCommand(cfg))

	if cfg.IsInsideApp() {
		scriptLoader := ac.ScriptLoader()

		err = scriptLoader.Load(filepath.Join(cfg.RootDir(), "cmd"))
		if err != nil {
			err = errors.Wrap(err, "failed to load user-defined commands")
		}

		udCmds := make([]*cobra.Command, 0)
		for _, name := range scriptLoader.Names() {
			udCmds = append(udCmds, newUserDefinedCommand(ac.UI(), scriptLoader, name))
		}
		if len(udCmds) > 0 {
			cmd.AddCommand(udCmds...)
		}
	}

	return cmd
}
