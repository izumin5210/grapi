package gencmd

type Option func(*Ctx)

func WithGenerateCommand(c *Command) Option {
	return func(ctx *Ctx) {
		ctx.GenerateCmd = c
	}
}

func WithDestroyCommand(c *Command) Option {
	return func(ctx *Ctx) {
		ctx.DestroyCmd = c
	}
}
