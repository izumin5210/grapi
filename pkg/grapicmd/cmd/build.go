package cmd

import (
	"fmt"
	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/srvc/appctx"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
)

func splitOptions(args []string) ([]string, []string) {
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

	scriptLoader := di.NewScriptLoader(ctx)
	errScriptLoad := scriptLoader.Load(ctx.RootDir.Join("cmd").String())

	isInsideApp := ctx.IsInsideApp()
	ui := di.NewUI(ctx)

	if !isInsideApp {
		fmt.Fprintln(os.Stderr, errors.New("protoc command should be execute inside a grapi application directory"))
		return nil
	}
	if errScriptLoad != nil {
		fmt.Fprintln(os.Stderr, errors.WithStack(errScriptLoad))
		return nil
	}

	return newBuildCommandMocked(scriptLoader, ui)
}

func newBuildCommandMocked(scriptLoader module.ScriptLoader, ui cli.UI) *cobra.Command {
	return &cobra.Command{
		Use:           "build [TARGET]... [-- BUILD_OPTIONS]",
		Short:         "Build commands",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(c *cobra.Command, args []string) error {
			names, opt := splitOptions(args)

			nameSet := make(map[string]bool, len(names))
			for _, n := range names {
				nameSet[n] = true
			}
			isAll := len(names) == 0

			ctx := appctx.Global()

			scriptLoaderNames := scriptLoader.Names()
			for _, name := range scriptLoaderNames {
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
