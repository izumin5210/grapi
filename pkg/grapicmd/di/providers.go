package di

import (
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

func ProvideUI(cfg *grapicmd.Config) clui.UI {
	return clui.New(cfg.OutWriter, cfg.InReader)
}

func ProvideCommandExecutor(cfg *grapicmd.Config, ui clui.UI) command.Executor {
	return command.NewExecutor(cfg.OutWriter, cfg.ErrWriter, cfg.InReader)
}

func ProvideGenerator(cfg *grapicmd.Config, ui clui.UI) module.Generator {
	return generator.New(
		cfg.Fs,
		ui,
		cfg.RootDir,
		cfg.ProtocConfig.ProtosDir,
		cfg.ProtocConfig.OutDir,
		cfg.ServerDir,
		cfg.Package,
		cfg.Version,
	)
}

func ProvideScriptLoader(cfg *grapicmd.Config, executor command.Executor) module.ScriptLoader {
	return script.NewLoader(cfg.Fs, executor, cfg.RootDir)
}

func ProvideGexConfig(cfg *grapicmd.Config) *gex.Config {
	return &gex.Config{
		OutWriter:  cfg.OutWriter,
		ErrWriter:  cfg.ErrWriter,
		InReader:   cfg.InReader,
		FS:         cfg.Fs,
		WorkingDir: cfg.RootDir,
		// TODO: set verbose flag
		// TODO: set logger
	}
}

func ProvideInitializeProjectUsecase(cfg *grapicmd.Config, gexCfg *gex.Config, ui clui.UI, generator module.Generator) usecase.InitializeProjectUsecase {
	return usecase.NewInitializeProjectUsecase(
		ui,
		generator,
		gexCfg,
		cfg.Version,
	)
}

func ProvideExecuteProtocUsecase(cfg *grapicmd.Config, gexCfg *gex.Config, ui clui.UI, executor command.Executor, generator module.Generator) usecase.ExecuteProtocUsecase {
	return usecase.NewExecuteProtocUsecase(
		cfg.ProtocConfig,
		cfg.Fs,
		ui,
		executor,
		gexCfg,
		cfg.RootDir,
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
