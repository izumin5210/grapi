// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package gencmd

import (
	cli "github.com/izumin5210/grapi/pkg/cli"
	grapicmd "github.com/izumin5210/grapi/pkg/grapicmd"
)

// Injectors from wire.go:

func NewApp(ctx *Ctx, command *Command) (*App, error) {
	ctx2 := ProvideGrapiCtx(ctx)
	fs := grapicmd.ProvideFS(ctx2)
	io := grapicmd.ProvideIO(ctx2)
	ui := cli.UIInstance(io)
	rootDir := grapicmd.ProvideRootDir(ctx2)
	fileSystem := ProvideTemplateFS(command)
	shouldRunFunc := ProvideShouldRun(command)
	generator := NewGenerator(fs, ui, rootDir, fileSystem, shouldRunFunc)
	app := &App{
		Generator: generator,
		UI:        ui,
	}
	return app, nil
}