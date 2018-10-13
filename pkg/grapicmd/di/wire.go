//+build wireinject

package di

import (
	"github.com/google/go-cloud/wire"

	"github.com/izumin5210/grapi/pkg/clui"
	"github.com/izumin5210/grapi/pkg/command"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
)

func NewUI(*grapicmd.Config) clui.UI {
	wire.Build(Set)
	return nil
}

func NewCommandExecutor(*grapicmd.Config) command.Executor {
	wire.Build(Set)
	return nil
}

func NewGenerator(*grapicmd.Config) module.Generator {
	wire.Build(Set)
	return nil
}

func NewScriptLoader(*grapicmd.Config) module.ScriptLoader {
	wire.Build(Set)
	return nil
}

func NewInitializeProjectUsecase(*grapicmd.Config) usecase.InitializeProjectUsecase {
	wire.Build(Set)
	return nil
}

func NewExecuteProtocUsecase(*grapicmd.Config) usecase.ExecuteProtocUsecase {
	wire.Build(Set)
	return nil
}
