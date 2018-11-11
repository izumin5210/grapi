// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package testing

import (
	cli "github.com/izumin5210/grapi/pkg/cli"
	gencmd "github.com/izumin5210/grapi/pkg/gencmd"
	grapicmd "github.com/izumin5210/grapi/pkg/grapicmd"
)

// Injectors from wire.go:

func NewTestApp(ctx *gencmd.Ctx, command *gencmd.Command, ui cli.UI) (*gencmd.App, error) {
	ctx2 := gencmd.ProvideGrapiCtx(ctx)
	fs := grapicmd.ProvideFS(ctx2)
	rootDir := grapicmd.ProvideRootDir(ctx2)
	fileSystem := gencmd.ProvideTemplateFS(command)
	shouldRunFunc := gencmd.ProvideShouldRun(command)
	generator := gencmd.NewGenerator(fs, ui, rootDir, fileSystem, shouldRunFunc)
	app := &gencmd.App{
		Generator: generator,
		UI:        ui,
	}
	return app, nil
}