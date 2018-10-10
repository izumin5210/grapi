package command

import "context"

// Executor is an interface for executing external commands.
type Executor interface {
	Exec(ctx context.Context, name string, opts ...Option) ([]byte, error)
}

// Command contains parameters for executing external commands.
type Command struct {
	Name        string
	Args        []string
	Dir         string
	Env         []string
	IOConnected bool
}

// Option specifies external command execution configurations.
type Option func(*Command)
