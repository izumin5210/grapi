package util

import (
	"path/filepath"
	"strings"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
	"github.com/pkg/errors"
	"github.com/serenize/snaker"
)

type ProtoParams struct {
	Proto struct {
		Path    string
		Package string
	}
	PbGo struct {
		Package    string
		ImportName string
	}
}

func BuildProtoParams(path string, rootDir cli.RootDir, protoOutDir string, pkg string) (out ProtoParams, err error) {
	if protoOutDir == "" {
		err = errors.New("protoOutDir is required")
		return
	}

	// github.com/foo/bar
	importPath, err := fs.GetImportPath(rootDir.String())
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// path => baz/qux/quux
	path = strings.Replace(snaker.CamelToSnake(path), "-", "_", -1)

	// baz/qux
	packagePath := filepath.Dir(path)

	// api/baz/qux
	pbgoPackagePath := filepath.Join(protoOutDir, packagePath)
	// qux_pb
	pbgoPackageName := filepath.Base(pbgoPackagePath) + "_pb"

	if packagePath == "." {
		pbgoPackagePath = protoOutDir
		pbgoPackageName = filepath.Base(pbgoPackagePath) + "_pb"
	}

	protoPackage := pkg
	if protoPackage == "" {
		protoPackageChunks := []string{}
		for _, pkg := range strings.Split(filepath.ToSlash(filepath.Join(importPath, protoOutDir)), "/") {
			chunks := strings.Split(pkg, ".")
			for i := len(chunks) - 1; i >= 0; i-- {
				protoPackageChunks = append(protoPackageChunks, chunks[i])
			}
		}
		// com.github.foo.bar.baz.qux
		protoPackage = strings.Join(protoPackageChunks, ".")
	}
	if dir := filepath.Dir(path); dir != "." {
		protoPackage = protoPackage + "." + strings.Replace(dir, string(filepath.Separator), ".", -1)
	}
	protoPackage = strings.Replace(protoPackage, "-", "_", -1)

	out.Proto.Path = path
	out.Proto.Package = strings.Replace(protoPackage, "-", "_", -1)
	out.PbGo.Package = filepath.ToSlash(filepath.Join(importPath, pbgoPackagePath))
	out.PbGo.ImportName = pbgoPackageName

	return
}
