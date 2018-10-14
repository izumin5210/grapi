package excmd

import "context"

// FakeExecutor do no-ops and records executed commands.
type FakeExecutor struct {
	Commands []*Command
}

// Exec implements Executor.Exec.
func (e *FakeExecutor) Exec(ctx context.Context, name string, opts ...Option) (_ []byte, _ error) {
	e.Commands = append(e.Commands, buildCommand(name, opts))
	return
}
