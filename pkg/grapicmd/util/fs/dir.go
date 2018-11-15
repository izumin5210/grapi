package fs

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

// CreateDirIfNotExists creates a directory if it does not exist.
func CreateDirIfNotExists(fs afero.Fs, path string) (err error) {
	err = fs.MkdirAll(path, 0755)
	zap.L().Debug("CreateDirIfNotExists", zap.String("path", path), zap.Error(err))
	return errors.Wrapf(err, "failed to create %q directory", path)
}
