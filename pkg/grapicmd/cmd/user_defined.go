package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

func newUserDefinedCommand(ui module.UI, scriptLoader module.ScriptLoader, name string) *cobra.Command {
	return &cobra.Command{
		Use:           name,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(c *cobra.Command, args []string) (err error) {
			script, ok := scriptLoader.Get(name)
			if !ok {
				err = errors.Wrapf(err, "faild to find subcommand %d", name)
			}

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
