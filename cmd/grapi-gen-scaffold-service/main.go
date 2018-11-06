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
		"scaffold-service",
		gencmd.WithGenerateCommand(NewGenerateCommand(svcgen.NewApp)),
	)
}

type CreateAppFunc func(*gencmd.Ctx, *gencmd.Command) (*svcgen.App, error)

func NewGenerateCommand(createApp CreateAppFunc) *gencmd.Command {
	var (
		skipTest bool
		resName  string
	)

	cmd := &gencmd.Command{
		Short:      "Generate a new service with standard methods",
		Args:       cobra.ExactArgs(1),
		TemplateFS: template.FS,
		BuildParams: func(c *gencmd.Command, args []string) (interface{}, error) {
			app, err := createApp(c.Ctx(), c)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			svcName := args[0]
			methods := []string{"list", "get", "create", "update", "delete"}

			params, err := app.ParamsBuilder.Build(svcName, resName, methods)
			return params, errors.WithStack(err)
		},
		PostRun: func(c *gencmd.Command, args []string) error {
			app, err := createApp(c.Ctx(), c)
			if err != nil {
				return errors.WithStack(err)
			}
			return errors.WithStack(app.ProtocWrapper.Exec(context.TODO()))
		},
	}

	cmd.Flags().BoolVarP(&skipTest, "skip-test", "T", false, "Skip test files")
	cmd.Flags().StringVar(&resName, "resource-name", "", "ResourceName to be used")

	return cmd
}
