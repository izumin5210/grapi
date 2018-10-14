package testingexcmd

import (
	"context"

	"github.com/izumin5210/grapi/pkg/excmd"
)

// FakeExecutor do no-ops and records executed commands.
type FakeExecutor struct {
	Commands []*excmd.Command
}

// Exec implements Executor.Exec.
func (e *FakeExecutor) Exec(ctx context.Context, name string, opts ...excmd.Option) (_ []byte, _ error) {
	e.Commands = append(e.Commands, excmd.BuildCommand(name, opts))
	return
}
