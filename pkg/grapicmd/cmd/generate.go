package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
)

func newGenerateCommand(cfg grapicmd.Config, ui module.UI, generatorFactory module.GeneratorFactory, commandFactory module.CommandFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate GENERATOR",
		Short:   "Generate new code",
		Aliases: []string{"g", "gen"},
	}

	cmd.AddCommand(newGenerateServiceCommand(cfg, ui, generatorFactory.Service(), commandFactory))
	cmd.AddCommand(newGenerateCommandCommand(cfg, ui, generatorFactory.Command()))

	return cmd
}

func newGenerateServiceCommand(cfg grapicmd.Config, ui module.UI, generator module.Generator, commandFactory module.CommandFactory) *cobra.Command {
	return &cobra.Command{
		Use:           "service NAME",
		Short:         "Generate a new service",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.IsInsideApp() {
				return errors.New("geneate command should execut inside a grapi applicaiton directory")
			}

			generateUsecase := usecase.NewGenerateServiceUsecase(ui, generator, cfg.RootDir())
			err := errors.WithStack(generateUsecase.Generate(args[0]))
			if err != nil {
				return err
			}

			protocUsecase := usecase.NewExecuteProtocUsecase(cfg.ProtocConfig(), cfg.Fs(), ui, commandFactory, cfg.RootDir())
			return errors.WithStack(protocUsecase.Perform())
		},
	}
}

func newGenerateCommandCommand(cfg grapicmd.Config, ui module.UI, generator module.Generator) *cobra.Command {
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
			return errors.WithStack(generator.Generate(cfg.RootDir(), data))
		},
	}
}
