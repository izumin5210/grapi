package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
)

func newDestroyCommand(cfg grapicmd.Config, ui module.UI, generatorFactory module.GeneratorFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "destroy GENERATOR",
		Short:   "Destroy codes",
		Aliases: []string{"d"},
	}

	cmd.AddCommand(newDestroyServiceCommand(cfg, ui, generatorFactory.Service()))
	cmd.AddCommand(newDestroyCommandCommand(cfg, ui, generatorFactory.Command()))

	return cmd
}

func newDestroyServiceCommand(cfg grapicmd.Config, ui module.UI, generator module.Generator) *cobra.Command {
	return &cobra.Command{
		Use:           "service NAME",
		Short:         "Destroy a service",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.IsInsideApp() {
				return errors.New("destroy command should execute inside a grapi applicaiton directory")
			}

			u := usecase.NewGenerateServiceUsecase(ui, generator, cfg.RootDir())
			return errors.WithStack(errors.WithStack(u.Destroy(args[0])))
		},
	}
}

func newDestroyCommandCommand(cfg grapicmd.Config, ui module.UI, generator module.Generator) *cobra.Command {
	return &cobra.Command{
		Use:   "command NAME",
		Short: "Destroy a command",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.IsInsideApp() {
				return errors.New("destroy command should execute inside a grapi applicaiton directory")
			}

			data := map[string]string{
				"name": args[0],
			}
			return errors.WithStack(generator.Destroy(cfg.RootDir(), data))
		},
	}
}
