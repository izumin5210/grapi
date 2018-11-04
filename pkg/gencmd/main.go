package gencmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Main(name string, opts ...Option) {
	var code int

	if err := run(name, opts...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		code = 1
	}

	os.Exit(code)
}

func run(name string, opts ...Option) error {
	ctx := defaultCtx()
	for _, f := range opts {
		f(ctx)
	}

	err := ctx.Init()
	if err != nil {
		return errors.Wrap(err, "failed to initialize context")
	}

	cmd := NewCommand(name, ctx)

	return errors.WithStack(cmd.Execute())
}

func NewCommand(name string, ctx *Ctx) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "grapi-gen-" + name,
	}

	setGenerateCommand(rootCmd, ctx)
	setDestroyCommand(rootCmd, ctx)

	return rootCmd
}

func setGenerateCommand(rootCmd *cobra.Command, ctx *Ctx) {
	cmd := ctx.GenerateCmd
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

func setDestroyCommand(rootCmd *cobra.Command, ctx *Ctx) {
	// TODO
}
