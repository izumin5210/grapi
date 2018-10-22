package cli

import (
	"io"
	"path/filepath"
)

type IO struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type RootDir string

func (d RootDir) String() string { return string(d) }

func (d RootDir) Join(elem ...string) string {
	return filepath.Join(append([]string{d.String()}, elem...)...)
}

func (d RootDir) BinDir() string {
	return d.Join("bin")
}
