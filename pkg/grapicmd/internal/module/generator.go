package module

// Generator creates files from templates and given params.
type Generator interface {
	ProjectGenerator
	ServiceGenerator
	CommandGenerator
}

// ProjectGenerator is an interface to build a new project.
type ProjectGenerator interface {
	GenerateProject(rootDir string, useHead bool) error
}

// ServiceGenerator is an interface to create or destroy gRPC services and implementations.
type ServiceGenerator interface {
	GenerateService(name string, methods ...string) error
	DestroyService(name string) error
}

// CommandGenerator is an interface to create or destroy user-defined command tempates.
type CommandGenerator interface {
	GenerateCommand(name string) error
	DestroyCommand(name string) error
}
