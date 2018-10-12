package cmd

import (
	grapicmd "github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type destroyCmd *cobra.Command
type destroySvcCmd *cobra.Command
type destroyCmdCmd *cobra.Command

func provideDestroyCommand(destroySvcCmd destroySvcCmd, destroyCmdCmd destroyCmdCmd) destroyCmd {
	cmd := &cobra.Command{
		Use:     "destroy GENERATOR",
		Short:   "Destroy codes",
		Aliases: []string{"d"},
	}

	cmd.AddCommand(destroySvcCmd)
	cmd.AddCommand(destroyCmdCmd)

	return cmd
}

func provideDestroyServiceCommand(cfg *grapicmd.Config, g module.Generator) destroySvcCmd {
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

			return errors.WithStack(g.DestroyService(args[0]))
		},
	}
}

func provideDestroyCommandCommand(cfg *grapicmd.Config, g module.Generator) destroyCmdCmd {
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

			return errors.WithStack(g.DestroyCommand(args[0]))
		},
	}
}
