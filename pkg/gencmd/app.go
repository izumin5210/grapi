package gencmd

import "github.com/izumin5210/grapi/pkg/cli"

type CreateAppFunc func(*Ctx, *Command) (*App, error)

type App struct {
	Generator Generator
	UI        cli.UI
}
