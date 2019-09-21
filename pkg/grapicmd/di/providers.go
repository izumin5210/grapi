package di

import (
	"github.com/google/wire"
	"github.com/izumin5210/gex"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/excmd"
	"github.com/izumin5210/grapi/pkg/gencmd"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/script"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
	"github.com/izumin5210/grapi/pkg/protoc"
)

func ProvideScriptLoader(ctx *grapicmd.Ctx, executor excmd.Executor) module.ScriptLoader {
	return script.NewLoader(ctx.FS, executor, ctx.RootDir.String())
}

func ProvideInitializeProjectUsecase(ctx *grapicmd.Ctx, gexCfg *gex.Config, ui cli.UI, fs afero.Fs, generator gencmd.Generator, excmd excmd.Executor) usecase.InitializeProjectUsecase {
	return usecase.NewInitializeProjectUsecase(
		ui,
		fs,
		generator,
		excmd,
		gexCfg,
	)
}

func ProvideShouldRun() gencmd.ShouldRunFunc { return nil }

var Set = wire.NewSet(
	grapicmd.CtxSet,
	protoc.WrapperSet,
	cli.UIInstance,
	excmd.NewExecutor,
	ProvideScriptLoader,
	gencmd.NewGenerator,
	fs.New,
	ProvideShouldRun,
	ProvideInitializeProjectUsecase,
)
