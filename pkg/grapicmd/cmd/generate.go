package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate/template"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
	"github.com/izumin5210/grapi/pkg/grapicmd/ui"
)

func newGenerateCommand(cfg grapicmd.Config, ui ui.UI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate GENERATOR",
		Short:   "Generate new code",
		Aliases: []string{"g", "gen"},
	}

	cmd.AddCommand(newGenerateServiceCommand(cfg, ui))
	cmd.AddCommand(newGenerateCommandCommand(cfg, ui))

	return cmd
}

func newGenerateServiceCommand(cfg grapicmd.Config, ui ui.UI) *cobra.Command {
	return &cobra.Command{
		Use:           "service NAME",
		Short:         "Generate a new service",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.IsInsideApp() {
				return errors.New("geneate command should execut inside a grapi applicaiton directory")
			}

			generator := generate.NewGenerator(cfg.Fs(), ui, cfg.RootDir())
			generateUsecase := usecase.NewGenerateServiceUsecase(ui, generator, cfg.RootDir())
			return errors.WithStack(generateUsecase.Perform(args[0]))
		},
	}
}

func newGenerateCommandCommand(cfg grapicmd.Config, ui ui.UI) *cobra.Command {
	return &cobra.Command{
		Use:   "command NAME",
		Short: "Generate a new command",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.IsInsideApp() {
				return errors.New("geneate command should execut inside a grapi applicaiton directory")
			}

			data := map[string]string{
				"name": args[0],
			}
			return generate.NewGenerator(cfg.Fs(), ui, cfg.RootDir()).Run(template.Command, data)
		},
	}
}
