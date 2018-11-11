package gencmd

import (
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/pkg/errors"
)

func defaultCtx() *Ctx {
	return &Ctx{
		Ctx: &grapicmd.Ctx{},
	}
}

// Ctx defines a context of a generator.
type Ctx struct {
	*grapicmd.Ctx

	CreateAppFunc CreateAppFunc
}

func (c *Ctx) apply(opts []Option) {
	for _, f := range opts {
		f(c)
	}
}

// CreateApp initializes dependencies.
func (c *Ctx) CreateApp(cmd *Command) (*App, error) {
	f := c.CreateAppFunc
	if c.CreateAppFunc == nil {
		f = NewApp
	}
	app, err := f(c, cmd)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return app, nil
}
