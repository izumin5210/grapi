package svcgen

import (
	"github.com/google/wire"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/protoc"
	"github.com/izumin5210/grapi/pkg/svcgen/params"
)

func ProvideParamsBuilder(rootDir cli.RootDir, protocCfg *protoc.Config, grapiCfg *grapicmd.Config) params.Builder {
	return params.NewBuilder(
		rootDir,
		protocCfg.ProtosDir,
		protocCfg.OutDir,
		grapiCfg.Grapi.ServerDir,
		grapiCfg.Package,
	)
}

var Set = wire.NewSet(
	ProvideParamsBuilder,
	App{},
)
