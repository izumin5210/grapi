package gencmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// New creates an Executor instnace.
func New(
	name string,
	generateCmd *Command,
	destroyCmd *Command,
	opts ...Option,
) Executor {
	ctx := defaultCtx()
	ctx.apply(opts)

	rootCmd := &cobra.Command{
		Use: "grapi-gen-" + name,
	}

	setGenerateCommand(rootCmd, generateCmd, ctx)
	setDestroyCommand(rootCmd, destroyCmd, ctx)

	return newExecutor(ctx, rootCmd)
}

func setGenerateCommand(rootCmd *cobra.Command, cmd *Command, ctx *Ctx) {
	if cmd == nil {
		return
	}

	ccmd := cmd.newCobraCommand()

	ccmd.RunE = func(_ *cobra.Command, args []string) error {
		app, err := ctx.CreateApp(cmd)
		if err != nil {
			return errors.WithStack(err)
		}

		params, err := cmd.BuildParams(cmd, args)
		if err != nil {
			return errors.WithStack(err)
		}

		err = app.Generator.Generate(params)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	}

	if ccmd.Use == "" {
		ccmd.Use = "generate"
	}

	cmd.ctx = ctx
	rootCmd.AddCommand(ccmd)
}

func setDestroyCommand(rootCmd *cobra.Command, cmd *Command, ctx *Ctx) {
	if cmd == nil {
		return
	}

	ccmd := cmd.newCobraCommand()

	ccmd.RunE = func(_ *cobra.Command, args []string) error {
		app, err := ctx.CreateApp(cmd)
		if err != nil {
			return errors.WithStack(err)
		}

		params, err := cmd.BuildParams(cmd, args)
		if err != nil {
			return errors.WithStack(err)
		}

		err = app.Generator.Destroy(params)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	}

	if ccmd.Use == "" {
		ccmd.Use = "destroy"
	}

	cmd.ctx = ctx
	rootCmd.AddCommand(ccmd)
}
