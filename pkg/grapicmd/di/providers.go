package di

import (
	"sync"

	"github.com/google/go-cloud/wire"
	"github.com/izumin5210/gex"
	"github.com/izumin5210/gex/pkg/tool"
	"github.com/pkg/errors"

	"github.com/izumin5210/grapi/pkg/clui"
	"github.com/izumin5210/grapi/pkg/excmd"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/script"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
	"github.com/izumin5210/grapi/pkg/protoc"
)

var (
	ui   clui.UI
	uiMu sync.Mutex

	gexCfg   *gex.Config
	gexCfgMu sync.Mutex

	toolRepo   tool.Repository
	toolRepoMu sync.Mutex
)

func ProvideUI(ctx *grapicmd.Ctx) clui.UI {
	uiMu.Lock()
	defer uiMu.Unlock()
	if ui == nil {
		ui = clui.New(ctx.OutWriter, ctx.InReader)
	}
	return ui
}

func ProvideCommandExecutor(ctx *grapicmd.Ctx, ui clui.UI) excmd.Executor {
	return excmd.NewExecutor(ctx.OutWriter, ctx.ErrWriter, ctx.InReader)
}

func ProvideGenerator(ctx *grapicmd.Ctx, ui clui.UI) module.Generator {
	return generator.New(
		ctx.FS,
		ui,
		ctx.RootDir,
		ctx.ProtocConfig.ProtosDir,
		ctx.ProtocConfig.OutDir,
		ctx.Config.Grapi.ServerDir,
		ctx.Config.Package,
		ctx.Version,
	)
}

func ProvideScriptLoader(ctx *grapicmd.Ctx, executor excmd.Executor) module.ScriptLoader {
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
			Execer:     ctx.Execer,
			WorkingDir: ctx.RootDir,
			// TODO: set verbose flag
			// TODO: set logger
		}
	}
	return gexCfg
}

func ProvideToolRepository(ctx *grapicmd.Ctx, gexCfg *gex.Config) (tool.Repository, error) {
	toolRepoMu.Lock()
	defer toolRepoMu.Unlock()
	if toolRepo == nil {
		var err error
		toolRepo, err = gexCfg.Create()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return toolRepo, nil
}

func ProvideProtocWrapper(ctx *grapicmd.Ctx, ui clui.UI, toolRepo tool.Repository) protoc.Wrapper {
	return protoc.NewWrapper(
		&ctx.ProtocConfig,
		ctx.FS,
		ctx.Execer,
		ui,
		toolRepo,
		ctx.RootDir,
		ctx.BinDir,
	)
}

func ProvideInitializeProjectUsecase(ctx *grapicmd.Ctx, gexCfg *gex.Config, ui clui.UI, generator module.Generator) usecase.InitializeProjectUsecase {
	return usecase.NewInitializeProjectUsecase(
		ui,
		generator,
		gexCfg,
		ctx.Version,
	)
}

var Set = wire.NewSet(
	ProvideUI,
	ProvideCommandExecutor,
	ProvideGenerator,
	ProvideScriptLoader,
	ProvideGexConfig,
	ProvideToolRepository,
	ProvideProtocWrapper,
	ProvideInitializeProjectUsecase,
)
