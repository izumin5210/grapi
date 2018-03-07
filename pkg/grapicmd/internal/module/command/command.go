package command

import (
	"bytes"
	"io"
	"os/exec"
	"sync"

	"github.com/izumin5210/clicontrib/clog"
	"github.com/pkg/errors"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

type command struct {
	name        string
	args        []string
	dir         string
	env         []string
	ioConnected bool
	outWriter   io.Writer
	errWriter   io.Writer
	inReader    io.Reader
}

func (c *command) SetDir(dir string) module.Command {
	c.dir = dir
	return c
}

func (c *command) AddEnv(key, value string) module.Command {
	c.env = append(c.env, key+"="+value)
	return c
}

func (c *command) ConnectIO() module.Command {
	c.ioConnected = true
	return c
}

func (c *command) Exec() ([]byte, error) {
	cmd := c.build()

	clog.Debug("execute", "command", cmd.Args, "dir", cmd.Dir)
	if !c.ioConnected {
		out, err := cmd.CombinedOutput()
		return out, errors.WithStack(err)
	}

	var wg sync.WaitGroup
	var buf bytes.Buffer

	closers := make([]func() error, 0, 2)
	wg.Add(2)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	go func() {
		defer wg.Done()
		io.Copy(c.outWriter, io.TeeReader(outReader, &buf))
	}()
	closers = append(closers, outReader.Close)

	go func() {
		defer wg.Done()
		io.Copy(c.errWriter, io.TeeReader(errReader, &buf))
	}()
	closers = append(closers, errReader.Close)

	cmd.Stdin = c.inReader

	err = cmd.Run()
	for _, c := range closers {
		c()
	}
	wg.Wait()

	return buf.Bytes(), err
}

func (c *command) build() *exec.Cmd {
	cmd := exec.Command(c.name, c.args...)
	cmd.Dir = c.dir
	cmd.Env = c.env
	return cmd
}
