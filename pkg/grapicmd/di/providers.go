package di

import (
	"github.com/google/wire"
	"github.com/izumin5210/gex"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/excmd"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/script"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
	"github.com/izumin5210/grapi/pkg/protoc"
)

func ProvideGenerator(ctx *grapicmd.Ctx, ui cli.UI) module.Generator {
	return generator.New(
		ctx.FS,
		ui,
	)
}

func ProvideScriptLoader(ctx *grapicmd.Ctx, executor excmd.Executor) module.ScriptLoader {
	return script.NewLoader(ctx.FS, executor, ctx.RootDir.String())
}

func ProvideInitializeProjectUsecase(ctx *grapicmd.Ctx, gexCfg *gex.Config, ui cli.UI, generator module.Generator, excmd excmd.Executor) usecase.InitializeProjectUsecase {
	return usecase.NewInitializeProjectUsecase(
		ui,
		generator,
		excmd,
		gexCfg,
	)
}

var Set = wire.NewSet(
	grapicmd.CtxSet,
	protoc.WrapperSet,
	cli.UIInstance,
	excmd.NewExecutor,
	ProvideGenerator,
	ProvideScriptLoader,
	ProvideInitializeProjectUsecase,
)
