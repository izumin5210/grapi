package gencmd

import (
	"github.com/google/wire"
	"github.com/rakyll/statik/fs"

	"github.com/izumin5210/clig/pkg/clib"
	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd"
)

func ProvideGrapiCtx(ctx *Ctx) *grapicmd.Ctx      { return ctx.Ctx }
func ProvideCtx(cmd *Command) *Ctx                { return cmd.Ctx() }
func ProvideShouldRun(cmd *Command) ShouldRunFunc { return cmd.ShouldRun }
func ProvidePath(root cli.RootDir) clib.Path      { return root.Path }

// Set contains providers for DI.
var Set = wire.NewSet(
	grapicmd.CtxSet,
	fs.New,
	ProvideGrapiCtx,
	ProvideCtx,
	ProvideShouldRun,
	ProvidePath,
	NewGenerator,
	App{},
)
