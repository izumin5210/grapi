package excmd

import (
	"os"
	"path/filepath"
	"strings"
)

func buildCommand(name string, opts []Option) *Command {
	c := &Command{Name: name, Env: os.Environ()}
	for _, f := range opts {
		f(c)
	}
	return c
}

// WithArgs sets arguments for a command.
func WithArgs(args ...string) Option {
	return func(c *Command) {
		c.Args = append(c.Args, args...)
	}
}

// WithDir sets a working directory for a command.
func WithDir(dir string) Option {
	return func(c *Command) {
		c.Dir = dir
	}
}

// WithPATH sets a PATH for a command.
func WithPATH(value string) Option {
	return func(c *Command) {
		for i := 0; i < len(c.Env); i++ {
			kv := strings.Split(c.Env[i], "=")
			if kv[0] == "PATH" {
				c.Env[i] = "PATH=" + value + string(filepath.ListSeparator) + kv[1]
				return
			}
		}
		WithEnv("PATH", value)(c)
	}
}

// WithEnv append a environment variable for a command.
func WithEnv(key, value string) Option {
	return func(c *Command) {
		c.Env = append(c.Env, key+"="+value)
	}
}

// WithIOConnected connects i/o with a command.
func WithIOConnected() Option {
	return func(c *Command) {
		c.IOConnected = true
	}
}
