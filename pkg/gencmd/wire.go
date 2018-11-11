//+build wireinject

package gencmd

import (
	"github.com/google/go-cloud/wire"
	"github.com/izumin5210/grapi/pkg/cli"
)

func newApp(*Ctx, *Command) (*App, error) {
	wire.Build(
		Set,
		cli.UIInstance,
	)
	return nil, nil
}
