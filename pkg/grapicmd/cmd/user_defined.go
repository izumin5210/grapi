package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

func newUserDefinedCommand(ui module.UI, script module.Script) *cobra.Command {
	return &cobra.Command{
		Use:           script.Name(),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(c *cobra.Command, args []string) (err error) {
			ui.Section(script.Name())
			ui.Subsection("Building...")
			err = errors.WithStack(script.Build())
			if err != nil {
				return
			}

			ui.Subsection("Starting...")
			return errors.WithStack(script.Run())
		},
	}
}
