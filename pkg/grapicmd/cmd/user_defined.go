package cmd

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/clui"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

func newUserDefinedCommands(cfg *grapicmd.Config) (cmds []*cobra.Command) {
	if !cfg.InsideApp {
		return
	}

	scriptLoader := di.NewScriptLoader(cfg)

	err := scriptLoader.Load(filepath.Join(cfg.RootDir, "cmd"))
	if err != nil {
		// TODO: log
		return
	}

	ui := di.NewUI(cfg)

	for _, name := range scriptLoader.Names() {
		cmds = append(cmds, newUserDefinedCommand(ui, scriptLoader, name))
	}

	return
}

func newUserDefinedCommand(ui clui.UI, scriptLoader module.ScriptLoader, name string) *cobra.Command {
	return &cobra.Command{
		Use:           name + " [-- BUILD_OPTIONS] [-- RUN_ARGS]",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(c *cobra.Command, args []string) (err error) {
			script, ok := scriptLoader.Get(name)
			if !ok {
				err = errors.Wrapf(err, "faild to find subcommand %s", name)
			}

			pos := len(args)
			for i, arg := range args {
				if arg == "--" {
					pos = i
					break
				}
			}
			var buildArgs, runArgs []string
			if pos == len(args) {
				buildArgs = args
			} else {
				buildArgs = args[:pos]
				runArgs = args[pos+1:]
			}

			ui.Section(script.Name())
			ui.Subsection("Building...")
			err = errors.WithStack(script.Build(buildArgs...))
			if err != nil {
				return
			}

			ui.Subsection("Starting...")
			return errors.WithStack(script.Run(runArgs...))
		},
	}
}
