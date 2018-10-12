package cmd

import (
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newDestroyCommand(ac di.AppComponent) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "destroy GENERATOR",
		Short:   "Destroy codes",
		Aliases: []string{"d"},
	}

	cmd.AddCommand(newDestroyServiceCommand(ac))
	cmd.AddCommand(newDestroyCommandCommand(ac))

	return cmd
}

func newDestroyServiceCommand(ac di.AppComponent) *cobra.Command {
	return &cobra.Command{
		Use:           "service NAME",
		Short:         "Destroy a service",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ac.Config().InsideApp {
				return errors.New("destroy command should execute inside a grapi application directory")
			}

			return errors.WithStack(errors.WithStack(ac.Generator().DestroyService(args[0])))
		},
	}
}

func newDestroyCommandCommand(ac di.AppComponent) *cobra.Command {
	return &cobra.Command{
		Use:           "command NAME",
		Short:         "Destroy a command",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ac.Config().InsideApp {
				return errors.New("destroy command should execute inside a grapi application directory")
			}

			return errors.WithStack(ac.Generator().DestroyCommand(args[0]))
		},
	}
}
