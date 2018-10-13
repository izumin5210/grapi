package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

func newGenerateCommand(cfg *grapicmd.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate GENERATOR",
		Short:   "Generate new code",
		Aliases: []string{"g", "gen"},
	}

	cmd.AddCommand(newGenerateServiceCommand(cfg))
	cmd.AddCommand(newGenerateScaffoldServiceCommand(cfg))
	cmd.AddCommand(newGenerateCommandCommand(cfg))

	return cmd
}

func newGenerateServiceCommand(cfg *grapicmd.Config) *cobra.Command {
	var (
		skipTest     bool
		resourceName string
	)

	cmd := &cobra.Command{
		Use:           "service NAME [METHODS...]",
		Short:         "Generate a new service",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.InsideApp {
				return errors.New("geneate command should execut inside a grapi application directory")
			}

			ui := di.NewUI(cfg)

			ui.Section("Generate service")
			genCfg := module.ServiceGenerationConfig{ResourceName: resourceName, Methods: args[1:], SkipTest: skipTest}
			err := errors.WithStack(di.NewGenerator(cfg).GenerateService(args[0], genCfg))
			if err != nil {
				return err
			}

			return errors.WithStack(di.NewExecuteProtocUsecase(cfg).Perform())
		},
	}

	cmd.PersistentFlags().BoolVarP(&skipTest, "skip-test", "T", false, "Skip test files")
	cmd.PersistentFlags().StringVar(&resourceName, "resource-name", "", "ResourceName to be used")

	return cmd
}

func newGenerateScaffoldServiceCommand(cfg *grapicmd.Config) *cobra.Command {
	var (
		skipTest     bool
		resourceName string
	)

	cmd := &cobra.Command{
		Use:           "scaffold-service NAME",
		Short:         "Generate a new service with standard methods",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.InsideApp {
				return errors.New("geneate command should execut inside a grapi application directory")
			}

			ui := di.NewUI(cfg)

			ui.Section("Scaffold service")
			genCfg := module.ServiceGenerationConfig{ResourceName: resourceName, SkipTest: skipTest}
			err := errors.WithStack(di.NewGenerator(cfg).ScaffoldService(args[0], genCfg))
			if err != nil {
				return err
			}

			return errors.WithStack(di.NewExecuteProtocUsecase(cfg).Perform())
		},
	}

	cmd.PersistentFlags().BoolVarP(&skipTest, "skip-test", "T", false, "Skip test files")
	cmd.PersistentFlags().StringVar(&resourceName, "resource-name", "", "ResourceName to be used")

	return cmd
}

func newGenerateCommandCommand(cfg *grapicmd.Config) *cobra.Command {
	return &cobra.Command{
		Use:           "command NAME",
		Short:         "Generate a new command",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.InsideApp {
				return errors.New("geneate command should execut inside a grapi application directory")
			}

			return errors.WithStack(di.NewGenerator(cfg).GenerateCommand(args[0]))
		},
	}
}
