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
	gencmd.Main(
		"service",
		gencmd.WithGenerateCommand(NewGenerateCommand(svcgen.NewApp)),
	)
}

type CreateAppFunc func(*gencmd.Ctx, *gencmd.Command) (*svcgen.App, error)

func NewGenerateCommand(createApp CreateAppFunc) *gencmd.Command {
	var (
		skipTest bool
		resName  string
		app      *svcgen.App
	)

	cmd := &gencmd.Command{
		Short:      "Generate a new service",
		Args:       cobra.MinimumNArgs(1),
		TemplateFS: template.FS,
		PreRun: func(c *gencmd.Command, args []string) error {
			var err error
			app, err = createApp(c.Ctx(), c)
			return errors.WithStack(err)
		},
		BuildParams: func(c *gencmd.Command, args []string) (interface{}, error) {
			svcName := args[0]
			methods := args[1:]

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
