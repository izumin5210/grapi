package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd/di"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
)

func newGenerateCommand(ac di.AppComponent) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate GENERATOR",
		Short:   "Generate new code",
		Aliases: []string{"g", "gen"},
	}

	cmd.AddCommand(newGenerateServiceCommand(ac))
	cmd.AddCommand(newGenerateScaffoldServiceCommand(ac))
	cmd.AddCommand(newGenerateCommandCommand(ac))

	return cmd
}

func newGenerateServiceCommand(ac di.AppComponent) *cobra.Command {
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
			cfg := ac.Config()

			if !cfg.IsInsideApp() {
				return errors.New("geneate command should execut inside a grapi application directory")
			}

			ac.UI().Section("Generate service")
			genCfg := module.ServiceGenerationConfig{ResourceName: resourceName, Methods: args[1:], SkipTest: skipTest}
			err := errors.WithStack(ac.Generator().GenerateService(args[0], genCfg))
			if err != nil {
				return err
			}

			protocUsecase := usecase.NewExecuteProtocUsecase(cfg.ProtocConfig(), cfg.Fs(), ac.UI(), ac.CommandExecutor(), ac.GexConfig(), cfg.RootDir())
			return errors.WithStack(protocUsecase.Perform())
		},
	}

	cmd.PersistentFlags().BoolVarP(&skipTest, "skip-test", "T", false, "Skip test files")
	cmd.PersistentFlags().StringVar(&resourceName, "resource-name", "", "ResourceName to be used")

	return cmd
}

func newGenerateScaffoldServiceCommand(ac di.AppComponent) *cobra.Command {
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
			cfg := ac.Config()

			if !cfg.IsInsideApp() {
				return errors.New("geneate command should execut inside a grapi application directory")
			}

			ac.UI().Section("Scaffold service")
			genCfg := module.ServiceGenerationConfig{ResourceName: resourceName, SkipTest: skipTest}
			err := errors.WithStack(ac.Generator().ScaffoldService(args[0], genCfg))
			if err != nil {
				return err
			}

			protocUsecase := usecase.NewExecuteProtocUsecase(cfg.ProtocConfig(), cfg.Fs(), ac.UI(), ac.CommandExecutor(), ac.GexConfig(), cfg.RootDir())
			return errors.WithStack(protocUsecase.Perform())
		},
	}

	cmd.PersistentFlags().BoolVarP(&skipTest, "skip-test", "T", false, "Skip test files")
	cmd.PersistentFlags().StringVar(&resourceName, "resource-name", "", "ResourceName to be used")

	return cmd
}

func newGenerateCommandCommand(ac di.AppComponent) *cobra.Command {
	return &cobra.Command{
		Use:           "command NAME",
		Short:         "Generate a new command",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ac.Config().IsInsideApp() {
				return errors.New("geneate command should execut inside a grapi application directory")
			}

			return errors.WithStack(ac.Generator().GenerateCommand(args[0]))
		},
	}
}
