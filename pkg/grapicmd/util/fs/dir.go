package fs

import (
	"github.com/izumin5210/clicontrib/pkg/clog"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// CreateDirIfNotExists creates a directory if it does not exist.
func CreateDirIfNotExists(fs afero.Fs, path string) (err error) {
	err = fs.MkdirAll(path, 0755)
	clog.Debug("CreateDirIfNotExists", "path", path, "error", err)
	return errors.Wrapf(err, "failed to create %q directory", path)
}
