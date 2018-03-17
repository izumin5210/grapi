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

func (s *script) Build(args ...string) error {
	err := fs.CreateDirIfNotExists(s.fs, filepath.Dir(s.binPath))
	if err != nil {
		return errors.WithStack(err)
	}

	cmd := s.commandFactory.Create(s.buildCmd(args))
	_, err = cmd.ConnectIO().SetDir(s.rootDir).Exec()
	if err != nil {
		return errors.Wrapf(err, "failed to build %v", s.srcPaths)
	}

	return nil
}

func (s *script) Run(args ...string) error {
	cmd := s.commandFactory.Create(append([]string{s.binPath}, args...))
	_, err := cmd.ConnectIO().SetDir(s.rootDir).Exec()
	return errors.WithStack(err)
}

func (s *script) buildCmd(args []string) []string {
	cmd := make([]string, 0, 3+len(args)+len(s.srcPaths))
	cmd = append(cmd, "go", "build", "-o="+s.binPath)
	cmd = append(cmd, args...)
	cmd = append(cmd, s.srcPaths...)
	return cmd
}
