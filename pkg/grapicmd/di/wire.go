//+build wireinject

package di

import (
	"github.com/google/go-cloud/wire"

	"github.com/izumin5210/grapi/pkg/clui"
	"github.com/izumin5210/grapi/pkg/excmd"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
	"github.com/izumin5210/grapi/pkg/protoc"
)

func NewUI(*grapicmd.Ctx) clui.UI {
	wire.Build(Set)
	return nil
}

func NewCommandExecutor(*grapicmd.Ctx) excmd.Executor {
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

func NewProtocWrapper(*grapicmd.Ctx) (protoc.Wrapper, error) {
	wire.Build(Set)
	return nil, nil
}

func NewInitializeProjectUsecase(*grapicmd.Ctx) usecase.InitializeProjectUsecase {
	wire.Build(Set)
	return nil
}
