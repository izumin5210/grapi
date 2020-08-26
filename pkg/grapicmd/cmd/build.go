package cmd

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/srvc/appctx"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
)

func getBuildOption(args []string) ([]string, []string) {
	pos := len(args)
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			pos = i
			break
		}
	}
	return args[:pos], args[pos:]
}

func newBuildCommand(ctx *grapicmd.Ctx) *cobra.Command {
	return &cobra.Command{
		Use:           "build [TARGET]... [-- BUILD_OPTIONS]",
		Short:         "Build commands",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(c *cobra.Command, args []string) error {
			if !ctx.IsInsideApp() {
				return errors.New("protoc command should be execute inside a grapi application directory")
			}

			arg, opt := getBuildOption(args)

			nameSet := make(map[string]bool, len(arg))
			for _, n := range arg {
				nameSet[n] = true
			}
			isAll := len(arg) == 0

			scriptLoader := di.NewScriptLoader(ctx)
			ui := di.NewUI(ctx)

			err := scriptLoader.Load(ctx.RootDir.Join("cmd").String())
			if err != nil {
				return errors.WithStack(err)
			}

			ctx := appctx.Global()

			for _, name := range scriptLoader.Names() {
				script, ok := scriptLoader.Get(name)
				if ok && (isAll || nameSet[script.Name()]) {
					ui.Subsection("Building " + script.Name())
					err := script.Build(ctx, opt...)
					if err != nil {
						return errors.WithStack(err)
					}
				}
			}

			return nil
		},
	}
}
