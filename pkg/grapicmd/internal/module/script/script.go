package script

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

type script struct {
	fs             afero.Fs
	commandFactory module.CommandFactory
	rootDir        string
	name, binPath  string
	srcPaths       []string
}

func (s *script) Name() string {
	return s.name
}

func (s *script) Build() error {
	err := fs.CreateDirIfNotExists(s.fs, filepath.Dir(s.binPath))
	if err != nil {
		return errors.WithStack(err)
	}

	cmd := s.commandFactory.Create(append([]string{"go", "build", "-v", "-o=" + s.binPath}, s.srcPaths...))
	_, err = cmd.ConnectIO().SetDir(s.rootDir).Exec()
	if err != nil {
		return errors.Wrapf(err, "failed to build %v", s.srcPaths)
	}

	return nil
}

func (s *script) Run() error {
	if ok, err := afero.Exists(s.fs, s.binPath); err != nil || !ok {
		err = s.Build()
		if err != nil {
			return errors.WithStack(err)
		}
	}
	cmd := s.commandFactory.Create([]string{s.binPath})
	_, err := cmd.ConnectIO().SetDir(s.rootDir).Exec()
	return errors.WithStack(err)
}
