package cmd

import (
	"path/filepath"

	"github.com/izumin5210/clicontrib/pkg/ccmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/command"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/script"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/ui"
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

	ui := ui.New(cfg.OutWriter(), cfg.InReader())
	generator := generator.New(
		cfg.Fs(),
		ui,
		cfg.RootDir(),
		cfg.ProtocConfig().ProtosDir,
		cfg.ProtocConfig().OutDir,
		cfg.ServerDir(),
		cfg.Version(),
	)
	commandFactory := command.NewFactory(cfg.OutWriter(), cfg.ErrWriter(), cfg.InReader())
	scriptLoader := script.NewLoader(cfg.Fs(), commandFactory, cfg.RootDir())

	cmd.AddCommand(newInitCommand(cfg, ui, generator, commandFactory))
	cmd.AddCommand(newGenerateCommand(cfg, ui, generator, commandFactory))
	cmd.AddCommand(newDestroyCommand(cfg, generator))
	cmd.AddCommand(newProtocCommand(cfg, ui, commandFactory))
	cmd.AddCommand(newBuildCommand(cfg, ui, scriptLoader))
	cmd.AddCommand(newVersionCommand(cfg))

	if cfg.IsInsideApp() {
		err = scriptLoader.Load(filepath.Join(cfg.RootDir(), "cmd"))
		if err != nil {
			err = errors.Wrap(err, "failed to load user-defined commands")
		}

		udCmds := make([]*cobra.Command, 0)
		for _, name := range scriptLoader.Names() {
			udCmds = append(udCmds, newUserDefinedCommand(ui, scriptLoader, name))
		}
		if len(udCmds) > 0 {
			cmd.AddCommand(udCmds...)
		}
	}

	return cmd
}
