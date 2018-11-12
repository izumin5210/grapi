package main

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/gencmd"
	"github.com/izumin5210/grapi/pkg/svcgen"
	"github.com/izumin5210/grapi/pkg/svcgen/template"
)

func main() {
	buildCommand(svcgen.NewApp).MustExecute()
}

func buildCommand(createAppFunc svcgen.CreateAppFunc, opts ...gencmd.Option) gencmd.Executor {
	return gencmd.New(
		"scaffold-service",
		newGenerateCommand(createAppFunc),
		newDestroyCommand(createAppFunc),
		opts...,
	)
}

func newGenerateCommand(createApp svcgen.CreateAppFunc) *gencmd.Command {
	var (
		skipTest bool
		resName  string
		app      *svcgen.App
	)

	cmd := &gencmd.Command{
		Use:             "generate NAME [flags]",
		Short:           "Generate a new service with standard methods",
		Args:            cobra.ExactArgs(1),
		TemplateFS:      template.FS,
		ShouldInsideApp: true,
		PreRun: func(c *gencmd.Command, args []string) error {
			var err error
			app, err = createApp(c)
			return errors.WithStack(err)
		},
		BuildParams: func(c *gencmd.Command, args []string) (interface{}, error) {
			svcName := args[0]
			methods := []string{"list", "get", "create", "update", "delete"}

			params, err := app.ParamsBuilder.Build(svcName, resName, methods)
			return params, errors.WithStack(err)
		},
		PostRun: func(c *gencmd.Command, args []string) error {
			return errors.WithStack(app.ProtocWrapper.Exec(context.TODO()))
		},
	}

	cmd.Flags().BoolVarP(&skipTest, "skip-test", "T", false, "Skip test files")
	cmd.Flags().StringVar(&resName, "resource-name", "", "ResourceName to be used")

	return cmd
}

func newDestroyCommand(createApp svcgen.CreateAppFunc) *gencmd.Command {
	var (
		app *svcgen.App
	)

	cmd := &gencmd.Command{
		Use:             "destroy NAME",
		Short:           "Destroy an existing service",
		Args:            cobra.ExactArgs(1),
		TemplateFS:      template.FS,
		ShouldInsideApp: true,
		PreRun: func(c *gencmd.Command, args []string) error {
			var err error
			app, err = createApp(c)
			return errors.WithStack(err)
		},
		BuildParams: func(c *gencmd.Command, args []string) (interface{}, error) {
			params, err := app.ParamsBuilder.Build(args[0], "", nil)
			return params, errors.WithStack(err)
		},
	}

	return cmd
}
