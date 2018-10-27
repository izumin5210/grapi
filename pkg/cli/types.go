package cli

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	colorable "github.com/mattn/go-colorable"
)

type IO struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func DefaultIO() *IO {
	io := &IO{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	if runtime.GOOS == "windows" {
		io.Out = colorable.NewColorableStdout()
		io.Err = colorable.NewColorableStderr()
	}
	return io
}

type RootDir string

func (d RootDir) String() string { return string(d) }

func (d RootDir) Join(elem ...string) string {
	return filepath.Join(append([]string{d.String()}, elem...)...)
}

func (d RootDir) BinDir() string {
	return d.Join("bin")
}
