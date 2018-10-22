package di

import (
	"sync"

	"github.com/google/go-cloud/wire"
	"github.com/izumin5210/gex"
	"github.com/izumin5210/gex/pkg/tool"
	"github.com/pkg/errors"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/excmd"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/script"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
	"github.com/izumin5210/grapi/pkg/protoc"
)

var (
	gexCfg   *gex.Config
	gexCfgMu sync.Mutex

	toolRepo   tool.Repository
	toolRepoMu sync.Mutex
)

func ProvideIO(ctx *grapicmd.Ctx) *cli.IO {
	return ctx.IO
}

func ProvideGenerator(ctx *grapicmd.Ctx, ui cli.UI) module.Generator {
	return generator.New(
		ctx.FS,
		ui,
		ctx.RootDir.String(),
		ctx.ProtocConfig.ProtosDir,
		ctx.ProtocConfig.OutDir,
		ctx.Config.Grapi.ServerDir,
		ctx.Config.Package,
		ctx.Build.Version,
	)
}

func ProvideScriptLoader(ctx *grapicmd.Ctx, executor excmd.Executor) module.ScriptLoader {
	return script.NewLoader(ctx.FS, executor, ctx.RootDir.String())
}

func ProvideGexConfig(ctx *grapicmd.Ctx) *gex.Config {
	gexCfgMu.Lock()
	defer gexCfgMu.Unlock()
	if gexCfg == nil {
		gexCfg = &gex.Config{
			OutWriter:  ctx.IO.Out,
			ErrWriter:  ctx.IO.Err,
			InReader:   ctx.IO.In,
			FS:         ctx.FS,
			Execer:     ctx.Execer,
			WorkingDir: ctx.RootDir.String(),
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

func ProvideProtocWrapper(ctx *grapicmd.Ctx, ui cli.UI, toolRepo tool.Repository) protoc.Wrapper {
	return protoc.NewWrapper(
		&ctx.ProtocConfig,
		ctx.FS,
		ctx.Execer,
		ui,
		toolRepo,
		ctx.RootDir,
	)
}

func ProvideInitializeProjectUsecase(ctx *grapicmd.Ctx, gexCfg *gex.Config, ui cli.UI, generator module.Generator) usecase.InitializeProjectUsecase {
	return usecase.NewInitializeProjectUsecase(
		ui,
		generator,
		gexCfg,
		ctx.Build.Version,
	)
}

var Set = wire.NewSet(
	ProvideIO,
	cli.UIInstance,
	excmd.NewExecutor,
	ProvideGenerator,
	ProvideScriptLoader,
	ProvideGexConfig,
	ProvideToolRepository,
	ProvideProtocWrapper,
	ProvideInitializeProjectUsecase,
)
