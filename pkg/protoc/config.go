package protoc

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Config stores setting params related protoc.
type Config struct {
	ImportDirs []string `mapstructure:"import_dirs"`
	ProtosDir  string   `mapstructure:"protos_dir"`
	OutDir     string   `mapstructure:"out_dir"`
	Plugins    []*Plugin
}

// ProtoFiles returns .proto file paths.
func (c *Config) ProtoFiles(fs afero.Fs, rootDir string) ([]string, error) {
	paths := []string{}
	protosDir := filepath.Join(rootDir, c.ProtosDir)
	err := afero.Walk(fs, protosDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		if !info.IsDir() && filepath.Ext(path) == ".proto" {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, errors.WithStack(err)
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
