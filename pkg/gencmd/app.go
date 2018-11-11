package gencmd

import "github.com/izumin5210/grapi/pkg/cli"

// CreateAppFunc initializes dependencies.
type CreateAppFunc func(*Ctx, *Command) (*App, error)

// App contains dependencies to execute a generator.
type App struct {
	Generator Generator
	UI        cli.UI
}
