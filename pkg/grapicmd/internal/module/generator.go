package module

// Generator creates files from templates and given params.
type Generator interface {
	ProjectGenerator
	ServiceGenerator
	CommandGenerator
}

// ProjectGenerator is an interface to build a new project.
type ProjectGenerator interface {
	GenerateProject(rootDir, pkgName string, useHead bool) error
}

// ServiceGenerator is an interface to create or destroy gRPC services and implementations.
type ServiceGenerator interface {
	GenerateService(name string, cfg ServiceGenerationConfig) error
	ScaffoldService(name string, cfg ServiceGenerationConfig) error
	DestroyService(name string) error
}

// ServiceGenerationConfig contains configurations for generating a new service.
type ServiceGenerationConfig struct {
	ResourceName string
	Methods      []string
	SkipTest     bool
}

// CommandGenerator is an interface to create or destroy user-defined command tempates.
type CommandGenerator interface {
	GenerateCommand(name string) error
	DestroyCommand(name string) error
}
