package script

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/izumin5210/clig/pkg/clib"
	"github.com/izumin5210/execx"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

type script struct {
	fs            afero.Fs
	io            *clib.IO
	exec          *execx.Executor
	rootDir       string
	name, binPath string
	srcPaths      []string
}

func (s *script) Name() string {
	return s.name
}

func (s *script) Build(args ...string) error {
	zap.L().Debug("build script", zap.String("name", s.name), zap.String("bin", s.binPath), zap.Strings("srcs", s.srcPaths))
	err := fs.CreateDirIfNotExists(s.fs, filepath.Dir(s.binPath))
	if err != nil {
		return errors.WithStack(err)
	}

	cmd := s.exec.CommandContext(context.TODO(), "go", s.buildArgs(args)...)
	cmd.Dir = s.rootDir
	cmd.Stdout = s.io.Out
	cmd.Stderr = s.io.Err
	cmd.Stdin = s.io.In
	err = cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to build %v", s.srcPaths)
	}

	return nil
}

func (s *script) Run(args ...string) error {
	cmd := s.exec.CommandContext(context.TODO(), s.binPath, args...)
	cmd.Dir = s.rootDir
	cmd.Stdout = s.io.Out
	cmd.Stderr = s.io.Err
	cmd.Stdin = s.io.In
	return errors.WithStack(cmd.Run())
}

func (s *script) buildArgs(args []string) []string {
	built := make([]string, 0, 3+len(args)+len(s.srcPaths))
	built = append(built, "build", "-o="+s.binPath)
	built = append(built, args...)
	built = append(built, s.srcPaths...)
	return built
}
