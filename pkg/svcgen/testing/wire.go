//+build wireinject

package testing

import (
	"github.com/google/go-cloud/wire"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/gencmd"
	"github.com/izumin5210/grapi/pkg/protoc"
	"github.com/izumin5210/grapi/pkg/svcgen"
)

func NewTestApp(*gencmd.Ctx, *gencmd.Command, protoc.Wrapper, cli.UI) (*svcgen.App, error) {
	wire.Build(
		gencmd.Set,
		svcgen.ProvideParamsBuilder,
		svcgen.App{},
	)
	return nil, nil
}
