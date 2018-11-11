package main

import (
	"github.com/izumin5210/grapi/cmd/grapi-gen-command/template"
	"github.com/izumin5210/grapi/pkg/gencmd"
	"github.com/spf13/cobra"
)

func main() {
	buildCommand().MustExecute()
}

func buildCommand(opts ...gencmd.Option) gencmd.Executor {
	return gencmd.New(
		"command",
		newGenerateCommand(),
		newDestroyCommand(),
		opts...,
	)
}

func newGenerateCommand() *gencmd.Command {
	return &gencmd.Command{
		Use:             "generate NAME",
		Short:           "Generate a new command",
		Args:            cobra.ExactArgs(1),
		TemplateFS:      template.FS,
		ShouldInsideApp: true,
		BuildParams: func(c *gencmd.Command, args []string) (interface{}, error) {
			return map[string]string{"name": args[0]}, nil
		},
	}
}

func newDestroyCommand() *gencmd.Command {
	return &gencmd.Command{
		Use:             "destroy NAME",
		Short:           "Destroy a existing command",
		Args:            cobra.ExactArgs(1),
		TemplateFS:      template.FS,
		ShouldInsideApp: true,
		BuildParams: func(c *gencmd.Command, args []string) (interface{}, error) {
			return map[string]string{"name": args[0]}, nil
		},
	}
}
