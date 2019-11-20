package module

import "context"

// Script represents an user-defined command.
type Script interface {
	Name() string
	Build(ctx context.Context, args ...string) error
	Run(ctx context.Context, args ...string) error
}

// ScriptLoader is a factory object for creating Script objects.
type ScriptLoader interface {
	Load(dir string) error
	Get(name string) (script Script, ok bool)
	Names() []string
}
