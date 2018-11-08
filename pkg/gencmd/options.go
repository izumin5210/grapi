package gencmd

import (
	"github.com/izumin5210/grapi/pkg/grapicmd"
)

type Option func(*Ctx)

func WithGrapiCtx(gctx *grapicmd.Ctx) Option {
	return func(ctx *Ctx) {
		ctx.Ctx = gctx
	}
}

func WithCreateAppFunc(f CreateAppFunc) Option {
	return func(ctx *Ctx) {
		ctx.CreateAppFunc = f
	}
}
