package command

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// Executor is an interface for executing commands.
type Executor interface {
	Exec(nameAndArgs []string, opts ...Option) ([]byte, error)
}

// NewExecutor creates a new Executor instance.
func NewExecutor(
	dir string,
	outWriter io.Writer,
	errWriter io.Writer,
	inReader io.Reader,
) Executor {
	return &executor{
		defaultDir: dir,
		outWriter:  outWriter,
		errWriter:  errWriter,
		inReader:   inReader,
	}
}

type executor struct {
	defaultDir string
	outWriter  io.Writer
	errWriter  io.Writer
	inReader   io.Reader
}

func (e *executor) Exec(nameAndArgs []string, opts ...Option) ([]byte, error) {
	c := e.createCommand(nameAndArgs, opts)
	out, err := c.Exec()
	return out, errors.WithStack(err)
}

func (e *executor) createCommand(nameAndArgs []string, opts []Option) *command {
	o := e.createOptions(opts)
	name := nameAndArgs[0]
	args := []string{}
	if len(nameAndArgs) > 1 {
		args = nameAndArgs[1:]
	}
	cmd := exec.Command(name, args...)
	cmd.Dir = o.Dir
	cmd.Env = o.Env
	return &command{
		cmd:  cmd,
		opts: o,
	}
}

func (e *executor) createOptions(opts []Option) *Options {
	o := &Options{
		Dir:       e.defaultDir,
		Env:       os.Environ(),
		OutWriter: e.outWriter,
		ErrWriter: e.errWriter,
		InReader:  e.inReader,
	}
	for _, f := range opts {
		f(o)
	}
	return o
}
