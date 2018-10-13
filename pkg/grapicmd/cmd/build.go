package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
)

func newBuildCommand(cfg *grapicmd.Config) *cobra.Command {
	return &cobra.Command{
		Use:           "build [TARGET]... [-- BUILD_OPTIONS]",
		Short:         "Build commands",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(c *cobra.Command, args []string) (err error) {
			if !cfg.InsideApp {
				return errors.New("protoc command should be execute inside a grapi application directory")
			}

			nameSet := make(map[string]bool, len(args))
			for _, n := range args {
				nameSet[n] = true
			}
			isAll := len(args) == 0

			scriptLoader := di.NewScriptLoader(cfg)
			ui := di.NewUI(cfg)

			for _, name := range scriptLoader.Names() {
				script, ok := scriptLoader.Get(name)
				if ok && (isAll || nameSet[script.Name()]) {
					ui.Subsection("Building " + script.Name())
					err := script.Build(args...)
					if err != nil {
						return errors.WithStack(err)
					}
				}
			}

			return nil
		},
	}
}
