package di

import (
	"sync"

	"github.com/google/go-cloud/wire"
	"github.com/izumin5210/gex"

	"github.com/izumin5210/grapi/pkg/clui"
	"github.com/izumin5210/grapi/pkg/command"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/script"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
)

var (
	ui   clui.UI
	uiMu sync.Mutex

	gexCfg   *gex.Config
	gexCfgMu sync.Mutex
)

func ProvideUI(ctx *grapicmd.Ctx) clui.UI {
	uiMu.Lock()
	defer uiMu.Unlock()
	if ui == nil {
		ui = clui.New(ctx.OutWriter, ctx.InReader)
	}
	return ui
}

func ProvideCommandExecutor(ctx *grapicmd.Ctx, ui clui.UI) command.Executor {
	return command.NewExecutor(ctx.OutWriter, ctx.ErrWriter, ctx.InReader)
}

func ProvideGenerator(ctx *grapicmd.Ctx, ui clui.UI) module.Generator {
	return generator.New(
		ctx.FS,
		ui,
		ctx.RootDir,
		ctx.ProtocConfig.ProtosDir,
		ctx.ProtocConfig.OutDir,
		ctx.Config.ServerDir,
		ctx.Config.Package,
		ctx.Version,
	)
}

func ProvideScriptLoader(ctx *grapicmd.Ctx, executor command.Executor) module.ScriptLoader {
	return script.NewLoader(ctx.FS, executor, ctx.RootDir)
}

func ProvideGexConfig(ctx *grapicmd.Ctx) *gex.Config {
	gexCfgMu.Lock()
	defer gexCfgMu.Unlock()
	if gexCfg == nil {
		gexCfg = &gex.Config{
			OutWriter:  ctx.OutWriter,
			ErrWriter:  ctx.ErrWriter,
			InReader:   ctx.InReader,
			FS:         ctx.FS,
			WorkingDir: ctx.RootDir,
			// TODO: set verbose flag
			// TODO: set logger
		}
	}
	return gexCfg
}

func ProvideInitializeProjectUsecase(ctx *grapicmd.Ctx, gexCfg *gex.Config, ui clui.UI, generator module.Generator) usecase.InitializeProjectUsecase {
	return usecase.NewInitializeProjectUsecase(
		ui,
		generator,
		gexCfg,
		ctx.Version,
	)
}

func ProvideExecuteProtocUsecase(ctx *grapicmd.Ctx, gexCfg *gex.Config, ui clui.UI, executor command.Executor, generator module.Generator) usecase.ExecuteProtocUsecase {
	return usecase.NewExecuteProtocUsecase(
		&ctx.ProtocConfig,
		ctx.FS,
		ui,
		executor,
		gexCfg,
		ctx.RootDir,
	)
}

var Set = wire.NewSet(
	ProvideUI,
	ProvideCommandExecutor,
	ProvideGenerator,
	ProvideScriptLoader,
	ProvideGexConfig,
	ProvideInitializeProjectUsecase,
	ProvideExecuteProtocUsecase,
)
