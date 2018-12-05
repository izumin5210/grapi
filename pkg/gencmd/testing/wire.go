//+build wireinject

package testing

import (
	"github.com/google/wire"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/gencmd"
)

func NewTestApp(*gencmd.Command, cli.UI) (*gencmd.App, error) {
	wire.Build(
		gencmd.Set,
	)
	return nil, nil
}
