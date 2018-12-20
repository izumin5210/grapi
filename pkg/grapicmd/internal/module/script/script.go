package script

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/izumin5210/grapi/pkg/excmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

type script struct {
	fs            afero.Fs
	excmd         excmd.Executor
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

	_, err = s.excmd.Exec(context.TODO(), "go", s.buildOpts(args)...)
	if err != nil {
		return errors.Wrapf(err, "failed to build %v", s.srcPaths)
	}

	return nil
}

func (s *script) Run(args ...string) error {
	_, err := s.excmd.Exec(
		context.TODO(),
		s.binPath,
		excmd.WithArgs(args...),
		excmd.WithDir(s.rootDir),
		excmd.WithIOConnected(),
	)
	return errors.WithStack(err)
}

func (s *script) buildOpts(args []string) []excmd.Option {
	built := make([]string, 0, 3+len(args)+len(s.srcPaths))
	built = append(built, "build", "-o="+s.binPath)
	built = append(built, args...)
	built = append(built, s.srcPaths...)
	return []excmd.Option{
		excmd.WithArgs(built...),
		excmd.WithDir(s.rootDir),
		excmd.WithIOConnected(),
	}
}
