package protoc

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Config stores setting params related protoc.
type Config struct {
	ImportDirs []string `mapstructure:"import_dirs"`
	ProtosDir  string   `mapstructure:"protos_dir"`
	OutDir     string   `mapstructure:"out_dir"`
	Plugins    []*Plugin
}

// OutDirOf returns a directory path of protoc result output path for given proto file.
func (c *Config) OutDirOf(rootDir string, protoPath string) (string, error) {
	protosDir := filepath.Join(rootDir, c.ProtosDir)
	relProtoDir, err := filepath.Rel(protosDir, filepath.Dir(protoPath))
	if strings.Contains(relProtoDir, "..") {
		return "", errors.Errorf(".proto files should be included in %s", c.ProtosDir)
	}
	if err != nil {
		return "", errors.Wrapf(err, ".proto files should be included in %s", c.ProtosDir)
	}

	return filepath.Join(c.OutDir, relProtoDir), nil
}

// Commands returns protoc command and arguments for given proto file.
func (c *Config) Commands(rootDir, protoPath string) ([][]string, error) {
	cmds := make([][]string, 0, len(c.Plugins))
	relProtoPath, _ := filepath.Rel(rootDir, protoPath)

	outDir, err := c.OutDirOf(rootDir, protoPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, p := range c.Plugins {
		args := []string{}
		args = append(args, "-I", filepath.Dir(relProtoPath))
		for _, dir := range c.ImportDirs {
			args = append(args, "-I", dir)
		}
		args = append(args, p.toProtocArg(outDir))
		args = append(args, relProtoPath)
		cmds = append(cmds, append([]string{"protoc"}, args...))
	}

	return cmds, nil
}
