package module

// GeneratorFactory is an interface for creating Generator instance.
type GeneratorFactory interface {
	Project() Generator
	Service() Generator
	Command() Generator
}

// Generator creates files from templates and given params.
type Generator interface {
	Generate(dir string, data interface{}) error
	Destroy(dir string, data interface{}) error
}
