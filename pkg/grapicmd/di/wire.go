//+build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/izumin5210/clig/pkg/clib"
	"github.com/izumin5210/gex/pkg/tool"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/excmd"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
	"github.com/izumin5210/grapi/pkg/protoc"
)

func NewUI(*grapicmd.Ctx) cli.UI {
	wire.Build(Set)
	return nil
}

func NewCommandExecutor(*grapicmd.Ctx) excmd.Executor {
	wire.Build(Set)
	return nil
}

func NewScriptLoader(*grapicmd.Ctx) module.ScriptLoader {
	wire.Build(Set)
	return nil
}

func NewToolRepository(*grapicmd.Ctx) (tool.Repository, error) {
	wire.Build(Set)
	return nil, nil
}

func NewProtocWrapper(*grapicmd.Ctx) (protoc.Wrapper, error) {
	wire.Build(Set)
	return nil, nil
}

func NewInitializeProjectUsecase(*grapicmd.Ctx, clib.Path) (usecase.InitializeProjectUsecase, error) {
	wire.Build(Set)
	return nil, nil
}
