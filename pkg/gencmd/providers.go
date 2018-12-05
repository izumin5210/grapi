package gencmd

import (
	"net/http"

	"github.com/google/wire"

	"github.com/izumin5210/grapi/pkg/grapicmd"
)

func ProvideGrapiCtx(ctx *Ctx) *grapicmd.Ctx         { return ctx.Ctx }
func ProvideCtx(cmd *Command) *Ctx                   { return cmd.Ctx() }
func ProvideTemplateFS(cmd *Command) http.FileSystem { return cmd.TemplateFS }
func ProvideShouldRun(cmd *Command) ShouldRunFunc    { return cmd.ShouldRun }

// Set contains providers for DI.
var Set = wire.NewSet(
	grapicmd.CtxSet,
	ProvideGrapiCtx,
	ProvideCtx,
	ProvideTemplateFS,
	ProvideShouldRun,
	NewGenerator,
	App{},
)
