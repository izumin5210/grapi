package cmd

import (
	grapicmd "github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newDestroyCommand(cfg *grapicmd.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "destroy GENERATOR",
		Short:   "Destroy codes",
		Aliases: []string{"d"},
	}

	cmd.AddCommand(newDestroyServiceCommand(cfg))
	cmd.AddCommand(newDestroyCommandCommand(cfg))

	return cmd
}

func newDestroyServiceCommand(cfg *grapicmd.Config) *cobra.Command {
	return &cobra.Command{
		Use:           "service NAME",
		Short:         "Destroy a service",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.InsideApp {
				return errors.New("destroy command should execute inside a grapi application directory")
			}

			return errors.WithStack(di.NewGenerator(cfg).DestroyService(args[0]))
		},
	}
}

func newDestroyCommandCommand(cfg *grapicmd.Config) *cobra.Command {
	return &cobra.Command{
		Use:           "command NAME",
		Short:         "Destroy a command",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.InsideApp {
				return errors.New("destroy command should execute inside a grapi application directory")
			}

			return errors.WithStack(di.NewGenerator(cfg).DestroyCommand(args[0]))
		},
	}
}
