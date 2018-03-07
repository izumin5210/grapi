package command

import (
	"io"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

// NewFactory creates a new module.CommandFactory instance.
func NewFactory(
	outWriter io.Writer,
	errWriter io.Writer,
	inReader io.Reader,
) module.CommandFactory {
	return &factory{
		outWriter: outWriter,
		errWriter: errWriter,
		inReader:  inReader,
	}
}

type factory struct {
	outWriter io.Writer
	errWriter io.Writer
	inReader  io.Reader
}

func (f *factory) Create(nameAndArgs []string) module.Command {
	name := nameAndArgs[0]
	args := make([]string, 0, len(nameAndArgs)-1)
	if len(nameAndArgs) > 1 {
		args = nameAndArgs[1:]
	}
	return &command{
		name:      name,
		args:      args,
		outWriter: f.outWriter,
		errWriter: f.errWriter,
		inReader:  f.inReader,
	}
}
