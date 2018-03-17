package module

// Script represents an user-defined command.
type Script interface {
	Name() string
	Build(args ...string) error
	Run(args ...string) error
}

// ScriptLoader is a factory object for creating Script objects.
type ScriptLoader interface {
	Load(dir string) error
	Get(name string) (script Script, ok bool)
	Names() []string
}
