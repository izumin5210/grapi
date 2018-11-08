package svcgen

import (
	"github.com/izumin5210/grapi/pkg/gencmd"
	"github.com/izumin5210/grapi/pkg/protoc"
	"github.com/izumin5210/grapi/pkg/svcgen/params"
)

type CreateAppFunc func(*gencmd.Ctx, *gencmd.Command) (*App, error)

type App struct {
	ProtocWrapper protoc.Wrapper
	ParamsBuilder params.Builder
}
