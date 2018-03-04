package command

import (
	"bytes"
	"io"
	"os/exec"
	"sync"

	"github.com/izumin5210/clicontrib/clog"
	"github.com/pkg/errors"
)

type command struct {
	cmd  *exec.Cmd
	opts *Options
}

func (c *command) Exec() ([]byte, error) {
	clog.Debug("execute", "command", c.cmd.Args, "dir", c.cmd.Dir)
	if !c.opts.IOConnected {
		out, err := c.cmd.CombinedOutput()
		return out, errors.WithStack(err)
	}

	var wg sync.WaitGroup
	var buf bytes.Buffer

	closers := make([]func() error, 0, 2)
	wg.Add(2)

	outReader, err := c.cmd.StdoutPipe()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	errReader, err := c.cmd.StderrPipe()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	go func() {
		defer wg.Done()
		io.Copy(c.opts.OutWriter, io.TeeReader(outReader, &buf))
	}()
	closers = append(closers, outReader.Close)

	go func() {
		defer wg.Done()
		io.Copy(c.opts.ErrWriter, io.TeeReader(errReader, &buf))
	}()
	closers = append(closers, errReader.Close)

	c.cmd.Stdin = c.opts.InReader

	err = c.cmd.Run()
	for _, c := range closers {
		c()
	}
	wg.Wait()

	return buf.Bytes(), err
}
