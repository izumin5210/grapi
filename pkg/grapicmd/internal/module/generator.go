package module

// Generator creates files from templates and given params.
type Generator interface {
	ProjectGenerator
}

// ProjectGenerator is an interface to build a new project.
type ProjectGenerator interface {
	GenerateProject(rootDir, pkgName string) error
}
