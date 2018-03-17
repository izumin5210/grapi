package command

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"os/signal"
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

func (c *command) Exec() (out []byte, err error) {
	var wg sync.WaitGroup

	cmd := c.build()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh)
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer recover()
		for sig := range sigCh {
			if !cmd.ProcessState.Exited() {
				cmd.Process.Signal(sig)
			}
		}
	}()

	clog.Debug("execute", "command", cmd.Args, "dir", cmd.Dir)
	if c.ioConnected {
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

		out = buf.Bytes()
	} else {
		out, err = cmd.CombinedOutput()
	}

	signal.Reset()
	close(sigCh)

	wg.Wait()
	return
}

func (c *command) build() *exec.Cmd {
	cmd := exec.Command(c.name, c.args...)
	cmd.Dir = c.dir
	cmd.Env = c.env
	return cmd
}
