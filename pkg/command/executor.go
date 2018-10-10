package command

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"

	"github.com/izumin5210/clicontrib/pkg/clog"
	"github.com/pkg/errors"
)

// NewExecutor creates a new Executor instance.
func NewExecutor(
	outWriter io.Writer,
	errWriter io.Writer,
	inReader io.Reader,
) Executor {
	return &executor{
		outWriter: outWriter,
		errWriter: errWriter,
		inReader:  inReader,
	}
}

type executor struct {
	outWriter io.Writer
	errWriter io.Writer
	inReader  io.Reader
}

func (e *executor) Exec(ctx context.Context, name string, opts ...Option) (out []byte, err error) {
	var wg sync.WaitGroup

	c := buildCommand(name, opts)
	clog.Debug("execute", "command", c)

	cmd := exec.CommandContext(ctx, c.Name, c.Args...)
	cmd.Dir = c.Dir
	cmd.Env = c.Env

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh)
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer recover()
		for sig := range sigCh {
			clog.Debug("signal received", "signal", sig)
			if cmd.ProcessState == nil || cmd.ProcessState.Exited() {
				break
			}
			cmd.Process.Signal(sig)
		}
	}()

	out, err = e.exec(c, cmd)
	if err != nil {
		err = errors.WithStack(err)
	}

	signal.Reset()
	close(sigCh)

	wg.Wait()
	return
}

func (e *executor) exec(c *Command, cmd *exec.Cmd) (out []byte, err error) {
	if c.IOConnected {
		var (
			buf bytes.Buffer
			wg  sync.WaitGroup
		)

		closers := make([]func() error, 0, 2)

		outReader, eerr := cmd.StdoutPipe()
		if eerr != nil {
			err = errors.WithStack(eerr)
			return
		}
		errReader, eerr := cmd.StderrPipe()
		if eerr != nil {
			err = errors.WithStack(eerr)
			return
		}

		wg.Add(2)
		go func() {
			defer wg.Done()
			io.Copy(e.outWriter, io.TeeReader(outReader, &buf))
		}()
		closers = append(closers, outReader.Close)
		go func() {
			defer wg.Done()
			io.Copy(e.errWriter, io.TeeReader(errReader, &buf))
		}()
		closers = append(closers, errReader.Close)

		cmd.Stdin = e.inReader

		err = cmd.Run()
		for _, c := range closers {
			c()
		}
		wg.Wait()

		out = buf.Bytes()
	} else {
		out, err = cmd.CombinedOutput()
	}

	return
}
