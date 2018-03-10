package module

// Script represents an user-defined command.
type Script interface {
	Name() string
	Build() error
	Run() error
}

// ScriptFactory is a factory object for creating Script objects.
type ScriptFactory interface {
	Create(entryFilePath string) Script
}
