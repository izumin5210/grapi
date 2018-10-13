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

func NewUI(*grapicmd.Ctx) clui.UI {
	wire.Build(Set)
	return nil
}

func NewCommandExecutor(*grapicmd.Ctx) command.Executor {
	wire.Build(Set)
	return nil
}

func NewGenerator(*grapicmd.Ctx) module.Generator {
	wire.Build(Set)
	return nil
}

func NewScriptLoader(*grapicmd.Ctx) module.ScriptLoader {
	wire.Build(Set)
	return nil
}

func NewInitializeProjectUsecase(*grapicmd.Ctx) usecase.InitializeProjectUsecase {
	wire.Build(Set)
	return nil
}

func NewExecuteProtocUsecase(*grapicmd.Ctx) usecase.ExecuteProtocUsecase {
	wire.Build(Set)
	return nil
}
