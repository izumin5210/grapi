//+build wireinject

package testing

import (
	"github.com/google/go-cloud/wire"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/gencmd"
)

func NewTestApp(*gencmd.Ctx, *gencmd.Command, cli.UI) (*gencmd.App, error) {
	wire.Build(
		gencmd.Set,
	)
	return nil, nil
}
