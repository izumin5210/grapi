package generator

import (
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

// New creates a module.Generator instance.
func New() module.Generator {
	return &generator{}
}

type generator struct {
	module.ProjectGenerator
	module.ServiceGenerator
	module.CommandGenerator
}
