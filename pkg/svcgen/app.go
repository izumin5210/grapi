package svcgen

import (
	"github.com/izumin5210/grapi/pkg/protoc"
	"github.com/izumin5210/grapi/pkg/svcgen/params"
)

type App struct {
	ProtocWrapper protoc.Wrapper
	ParamsBuilder params.Builder
}
