package cli

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	colorable "github.com/mattn/go-colorable"
)

// IO contains an input reader, an output writer and an error writer.
type IO struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

// DefaultIO returns a standard IO object.
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

// RootDir represents a project root directory.
type RootDir string

func (d RootDir) String() string { return string(d) }

// Join joins path elements to the root directory.
func (d RootDir) Join(elem ...string) string {
	return filepath.Join(append([]string{d.String()}, elem...)...)
}

// BinDir returns the directory path contains executable binaries.
func (d RootDir) BinDir() string {
	return d.Join("bin")
}
