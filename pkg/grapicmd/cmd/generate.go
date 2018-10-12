package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/clui"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
)

type generateCmd *cobra.Command
type genSvcCmd *cobra.Command
type genScaffoldSvcCmd *cobra.Command
type genCmdCmd *cobra.Command

func provideGenerateCommand(
	genSvcCmd genSvcCmd,
	genScaffoldSvcCmd genScaffoldSvcCmd,
	genCmdCmd genCmdCmd,
) generateCmd {
	cmd := &cobra.Command{
		Use:     "generate GENERATOR",
		Short:   "Generate new code",
		Aliases: []string{"g", "gen"},
	}

	cmd.AddCommand(genSvcCmd)
	cmd.AddCommand(genScaffoldSvcCmd)
	cmd.AddCommand(genCmdCmd)

	return cmd
}

func provideGenerateServiceCommand(cfg *grapicmd.Config, ui clui.UI, g module.Generator, u usecase.ExecuteProtocUsecase) genSvcCmd {
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

			ui.Section("Generate service")
			genCfg := module.ServiceGenerationConfig{ResourceName: resourceName, Methods: args[1:], SkipTest: skipTest}
			err := errors.WithStack(g.GenerateService(args[0], genCfg))
			if err != nil {
				return err
			}

			return errors.WithStack(u.Perform())
		},
	}

	cmd.PersistentFlags().BoolVarP(&skipTest, "skip-test", "T", false, "Skip test files")
	cmd.PersistentFlags().StringVar(&resourceName, "resource-name", "", "ResourceName to be used")

	return cmd
}

func provideGenerateScaffoldServiceCommand(cfg *grapicmd.Config, ui clui.UI, g module.Generator, u usecase.ExecuteProtocUsecase) genScaffoldSvcCmd {
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

			ui.Section("Scaffold service")
			genCfg := module.ServiceGenerationConfig{ResourceName: resourceName, SkipTest: skipTest}
			err := errors.WithStack(g.ScaffoldService(args[0], genCfg))
			if err != nil {
				return err
			}

			return errors.WithStack(u.Perform())
		},
	}

	cmd.PersistentFlags().BoolVarP(&skipTest, "skip-test", "T", false, "Skip test files")
	cmd.PersistentFlags().StringVar(&resourceName, "resource-name", "", "ResourceName to be used")

	return cmd
}

func provideGenerateCommandCommand(cfg *grapicmd.Config, g module.Generator) genCmdCmd {
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

			return errors.WithStack(g.GenerateCommand(args[0]))
		},
	}
}
