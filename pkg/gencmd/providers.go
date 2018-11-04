package gencmd

import (
	"net/http"

	"github.com/google/go-cloud/wire"

	"github.com/izumin5210/grapi/pkg/grapicmd"
)

func ProvideGrapiCtx(ctx *Ctx) *grapicmd.Ctx         { return ctx.Ctx }
func ProvideTemplateFS(cmd *Command) http.FileSystem { return cmd.TemplateFS }
func ProvideShouldRun(cmd *Command) ShouldRunFunc    { return cmd.ShouldRun }

var Set = wire.NewSet(
	grapicmd.CtxSet,
	ProvideGrapiCtx,
	ProvideTemplateFS,
	ProvideShouldRun,
	NewGenerator,
	App{},
)
