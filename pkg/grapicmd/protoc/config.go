package grapicmd

import (
	"path/filepath"

	"github.com/pkg/errors"
)

type Config struct {
	ImportDirs []string `mapstructure:"import_dirs"`
	ProtosDir  string   `mapstructure:"protos_dir"`
	OutDir     string   `mapstructure:"out_dir"`
	Plugins    []*Plugin
}

func (c *Config) OutDirOf(rootDir string, protoPath string) (string, error) {
	protosDir := filepath.Join(rootDir, c.ProtosDir)
	relProtoDir, err := filepath.Rel(protosDir, filepath.Dir(protoPath))
	if err != nil {
		return "", errors.Wrapf(err, ".proto files should be included in %s", c.ProtosDir)
	}

	return filepath.Join(c.OutDir, relProtoDir), nil
}

func (c *Config) Commands(rootDir, protoPath string) ([][]string, error) {
	cmds := make([][]string, 0, len(c.Plugins))
	relOutDir, err := c.OutDirOf(rootDir, protoPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	outDir := filepath.Join(rootDir, relOutDir)

	for _, p := range c.Plugins {
		args := []string{}
		args = append(args, "-I", filepath.Dir(protoPath))
		for _, dir := range c.ImportDirs {
			absDir := dir
			if !filepath.IsAbs(absDir) {
				absDir = filepath.Join(rootDir, absDir)
			}
			args = append(args, "-I", absDir)
		}
		args = append(args, p.toProtocArg(outDir))
		args = append(args, protoPath)
		cmds = append(cmds, append([]string{"protoc"}, args...))
	}

	return cmds, nil
}
