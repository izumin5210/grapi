package command

import (
	"io"
	"os"
)

// NewFactory creates a new Factory instance.
func NewFactory(
	outWriter io.Writer,
	errWriter io.Writer,
	inReader io.Reader,
) Factory {
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

func (f *factory) Create(nameAndArgs []string) Command {
	name := nameAndArgs[0]
	args := make([]string, 0, len(nameAndArgs)-1)
	if len(nameAndArgs) > 1 {
		args = nameAndArgs[1:]
	}
	return &command{
		name:      name,
		args:      args,
		env:       os.Environ(),
		outWriter: f.outWriter,
		errWriter: f.errWriter,
		inReader:  f.inReader,
	}
}
