// +build wireinject

package di

import (
	"github.com/google/go-cloud/wire"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

func initializeGenerator(cfg *grapicmd.Config) module.Generator {
	wire.Build(Set)
	return nil
}
