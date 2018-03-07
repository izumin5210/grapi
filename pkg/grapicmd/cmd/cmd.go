package cmd

import (
	"path/filepath"

	"github.com/izumin5210/clicontrib"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/command"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/ui"
)

// NewGrapiCommand creates a new command object.
func NewGrapiCommand(cfg grapicmd.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:           cfg.AppName(),
		Short:         "JSON API framework implemented with gRPC and Gateway",
		Long:          "",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	var cfgFile string
	cobra.OnInitialize(func() { cfg.Init(cfgFile) })
	clicontrib.HandleLogFlags(cmd)

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "./"+cfg.AppName()+".toml", "config file")
	commandFactory := command.NewFactory(cfg.OutWriter(), cfg.ErrWriter(), cfg.InReader())
	scriptFactory := internal.NewScriptFactory(cfg.Fs(), commandFactory, cfg.RootDir())

	ui := ui.New(cfg.OutWriter(), cfg.InReader())

	cmd.AddCommand(newInitCommand(cfg, ui, commandFactory))
	cmd.AddCommand(newGenerateCommand(cfg, ui, commandFactory))
	cmd.AddCommand(newProtocCommand(cfg, ui, commandFactory))
	cmd.AddCommand(newBuildCommand(cfg, ui, scriptFactory))
	cmd.AddCommand(newVersionCommand(cfg))

	udCmds := make([]*cobra.Command, 0)
	if cfg.IsInsideApp() {
		paths, err := afero.Glob(cfg.Fs(), filepath.Join(cfg.RootDir(), "cmd/*/run.go"))
		if err == nil {
			for _, path := range paths {
				udCmds = append(udCmds, newUserDefinedCommand(ui, scriptFactory.Create(path)))
			}
		}
	}
	if len(udCmds) > 0 {
		cmd.AddCommand(udCmds...)
	}

	return cmd
}
