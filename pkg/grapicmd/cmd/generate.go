package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate/template"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
)

func newGenerateCommand(cfg grapicmd.Config, ui module.UI, commandFactory module.CommandFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate GENERATOR",
		Short:   "Generate new code",
		Aliases: []string{"g", "gen"},
	}

	cmd.AddCommand(newGenerateServiceCommand(cfg, ui, commandFactory))
	cmd.AddCommand(newGenerateCommandCommand(cfg, ui))

	return cmd
}

func newGenerateServiceCommand(cfg grapicmd.Config, ui module.UI, commandFactory module.CommandFactory) *cobra.Command {
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
			err := errors.WithStack(generateUsecase.Perform(args[0]))
			if err != nil {
				return err
			}

			protocUsecase := usecase.NewExecuteProtocUsecase(cfg.ProtocConfig(), cfg.Fs(), ui, commandFactory, cfg.RootDir())
			return errors.WithStack(protocUsecase.Perform())
		},
	}
}

func newGenerateCommandCommand(cfg grapicmd.Config, ui module.UI) *cobra.Command {
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
