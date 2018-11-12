package main

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/serenize/snaker"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/cmd/grapi-gen-type/di"
	"github.com/izumin5210/grapi/cmd/grapi-gen-type/template"
	"github.com/izumin5210/grapi/pkg/gencmd"
	gencmdutil "github.com/izumin5210/grapi/pkg/gencmd/util"
	"github.com/izumin5210/grapi/pkg/grapicmd"
)

func main() {
	buildCommand(di.NewApp).MustExecute()
}

func buildCommand(createApp di.CreateAppFunc, opts ...gencmd.Option) gencmd.Executor {
	return gencmd.New(
		"type",
		newGenerateCommand(createApp),
		newDestroyCommand(),
		opts...,
	)
}

func newGenerateCommand(createApp di.CreateAppFunc) *gencmd.Command {
	var (
		app *di.App
	)

	return &gencmd.Command{
		Use:             "generate NAME",
		Short:           "Generate a new type",
		Args:            cobra.ExactArgs(1),
		TemplateFS:      template.FS,
		ShouldInsideApp: true,
		PreRun: func(c *gencmd.Command, args []string) error {
			var err error
			app, err = createApp(c)
			return errors.WithStack(err)
		},
		BuildParams: func(c *gencmd.Command, args []string) (interface{}, error) {
			return buildParams(args[0], c.Ctx().Ctx)
		},
		PostRun: func(c *gencmd.Command, args []string) error {
			return errors.WithStack(app.Protoc.Exec(context.TODO()))
		},
	}
}

func newDestroyCommand() *gencmd.Command {
	return &gencmd.Command{
		Use:             "destroy NAME",
		Short:           "Destroy a existing type",
		Args:            cobra.ExactArgs(1),
		TemplateFS:      template.FS,
		ShouldInsideApp: true,
		BuildParams: func(c *gencmd.Command, args []string) (interface{}, error) {
			return buildParams(args[0], c.Ctx().Ctx)
		},
	}
}

type params struct {
	Proto struct {
		Package string
	}
	PbGo struct {
		PackagePath string
		PackageName string
	}
	Name string
	Path string
}

func buildParams(path string, ctx *grapicmd.Ctx) (*params, error) {
	protoOutDir := ctx.ProtocConfig.OutDir
	if protoOutDir == "" {
		protoOutDir = filepath.Join("api")
	}
	protoOutDir = filepath.Join(protoOutDir, "type")

	protoParams, err := gencmdutil.BuildProtoParams(path, ctx.RootDir, protoOutDir, ctx.Config.Package)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// path => baz/qux/quux
	path = protoParams.Proto.Path

	params := new(params)
	params.Proto.Package = protoParams.Proto.Package
	params.PbGo.PackageName = protoParams.PbGo.ImportName
	params.PbGo.PackagePath = protoParams.PbGo.Package
	params.Name = snaker.SnakeToCamel(filepath.Base(path))
	params.Path = path
	return params, nil
}
