package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

func newDestroyCommand(cfg grapicmd.Config, generator module.Generator) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "destroy GENERATOR",
		Short:   "Destroy codes",
		Aliases: []string{"d"},
	}

	cmd.AddCommand(newDestroyServiceCommand(cfg, generator))
	cmd.AddCommand(newDestroyCommandCommand(cfg, generator))

	return cmd
}

func newDestroyServiceCommand(cfg grapicmd.Config, generator module.ServiceGenerator) *cobra.Command {
	return &cobra.Command{
		Use:           "service NAME",
		Short:         "Destroy a service",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.IsInsideApp() {
				return errors.New("destroy command should execute inside a grapi application directory")
			}

			return errors.WithStack(errors.WithStack(generator.DestroyService(args[0])))
		},
	}
}

func newDestroyCommandCommand(cfg grapicmd.Config, generator module.CommandGenerator) *cobra.Command {
	return &cobra.Command{
		Use:           "command NAME",
		Short:         "Destroy a command",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.IsInsideApp() {
				return errors.New("destroy command should execute inside a grapi application directory")
			}

			return errors.WithStack(generator.DestroyCommand(args[0]))
		},
	}
}
