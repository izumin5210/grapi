package gencmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Executor interface {
	Command() *cobra.Command
	Execute() error
	MustExecute()
}

func newExecutor(ctx *Ctx, cmd *cobra.Command) Executor {
	return &executorImpl{ctx: ctx, cmd: cmd}
}

type executorImpl struct {
	ctx *Ctx
	cmd *cobra.Command
}

func (c *executorImpl) Command() *cobra.Command {
	return c.cmd
}

func (c *executorImpl) Execute() error {
	err := c.ctx.Init()
	if err != nil {
		return errors.Wrap(err, "failed to initialize context")
	}

	return errors.WithStack(c.cmd.Execute())
}

func (c *executorImpl) MustExecute() {
	var code int

	if err := c.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		code = 1
	}

	os.Exit(code)
}
