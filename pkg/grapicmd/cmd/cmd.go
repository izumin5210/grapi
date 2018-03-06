package cmd

import (
	"path/filepath"

	"github.com/izumin5210/clicontrib"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/command"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal"
	"github.com/izumin5210/grapi/pkg/grapicmd/ui"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
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
	rootDir, inProj := fs.LookupRoot(cfg.Fs(), cfg.CurrentDir())
	scriptFactory := internal.NewScriptFactory(
		cfg.Fs(),
		command.NewExecutor(rootDir, cfg.OutWriter(), cfg.ErrWriter(), cfg.InReader()),
		rootDir,
	)

	ui := ui.New(cfg.OutWriter(), cfg.InReader())

	cmd.AddCommand(newInitCommand(cfg, ui))
	cmd.AddCommand(newGenerateCommand(cfg, ui))
	cmd.AddCommand(newProtocCommand(cfg, ui))
	cmd.AddCommand(newBuildCommand(cfg, ui, scriptFactory))
	cmd.AddCommand(newVersionCommand(cfg))

	udCmds := make([]*cobra.Command, 0)
	if inProj {
		paths, err := afero.Glob(cfg.Fs(), filepath.Join(rootDir, "cmd/*/run.go"))
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
