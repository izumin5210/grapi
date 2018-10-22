package cli

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
)

var (
	requiredFiles = []string{"grapi.toml"}
	requiredDirs  = []string{"app", "api", "cmd"}
)

// LookupRoot returns the application's root directory if the current directory is inside of a grapi application.
func LookupRoot(fs afero.Fs, dir string) (RootDir, bool) {
	var eg errgroup.Group

	for _, f := range requiredFiles {
		f := f
		eg.Go(func() error {
			ok, err := afero.Exists(fs, filepath.Join(dir, f))
			if err != nil || !ok {
				return errors.Errorf("%s does not exist", f)
			}
			return nil
		})
	}

	for _, d := range requiredDirs {
		d := d
		eg.Go(func() error {
			ok, err := afero.DirExists(fs, filepath.Join(dir, d))
			if err != nil || !ok {
				return errors.Errorf("%s does not exist", d)
			}
			return nil
		})
	}

	if eg.Wait() == nil {
		return RootDir(dir), true
	}

	p := dir[len(filepath.VolumeName(dir)):]
	if p == string(filepath.Separator) {
		return "", false
	}

	return LookupRoot(fs, filepath.Dir(dir))
}
