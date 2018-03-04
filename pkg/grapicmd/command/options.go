package command

import "io"

// Option is an option for executing commands.
type Option func(*Options)

// Options contains configurations for executing commands.
type Options struct {
	Dir         string
	Env         []string
	IOConnected bool
	OutWriter   io.Writer
	ErrWriter   io.Writer
	InReader    io.Reader
}

// WithDir returns an Option that sets the working directory of commands.
func WithDir(dir string) Option {
	return func(o *Options) {
		o.Dir = dir
	}
}

// WithEnv returns an Option that add an environment variable.
func WithEnv(key, value string) Option {
	return func(o *Options) {
		o.Env = append(o.Env, key+"="+value)
	}
}

// WithIOConnected returns an Option that connects with command's stdio.
func WithIOConnected() Option {
	return func(o *Options) {
		o.IOConnected = true
	}
}
