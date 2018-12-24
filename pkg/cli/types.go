package cli

import (
	"github.com/izumin5210/clig/pkg/clib"
)

// RootDir represents a project root directory.
type RootDir struct {
	clib.Path
}

// BinDir returns the directory path contains executable binaries.
func (d *RootDir) BinDir() clib.Path {
	return d.Join("bin")
}
