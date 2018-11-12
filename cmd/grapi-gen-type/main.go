package main

import (
	"path/filepath"
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
	"github.com/serenize/snaker"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/cmd/grapi-gen-type/template"
	"github.com/izumin5210/grapi/pkg/gencmd"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

func main() {
	buildCommand().MustExecute()
}

func buildCommand(opts ...gencmd.Option) gencmd.Executor {
	return gencmd.New(
		"type",
		newGenerateCommand(),
		newDestroyCommand(),
		opts...,
	)
}

func newGenerateCommand() *gencmd.Command {
	return &gencmd.Command{
		Use:             "generate NAME",
		Short:           "Generate a new type",
		Args:            cobra.ExactArgs(1),
		TemplateFS:      template.FS,
		ShouldInsideApp: true,
		BuildParams: func(c *gencmd.Command, args []string) (interface{}, error) {
			return buildParams(args[0], c.Ctx().Ctx)
		},
	}
}

func newDestroyCommand() *gencmd.Command {
	return &gencmd.Command{
		Use:             "destroy NAME",
		Short:           "Destroy a existing type",
		Args:            cobra.ExactArgs(1),
		TemplateFS:      template.FS,
		ShouldInsideApp: true,
		BuildParams: func(c *gencmd.Command, args []string) (interface{}, error) {
			return buildParams(args[0], c.Ctx().Ctx)
		},
	}
}

type params struct {
	Proto struct {
		Package string
	}
	PbGo struct {
		PackagePath string
		PackageName string
	}
	Name string
	Path string
}

func buildParams(path string, ctx *grapicmd.Ctx) (*params, error) {
	protoOutDir := ctx.ProtocConfig.OutDir
	if protoOutDir == "" {
		protoOutDir = filepath.Join("api")
	}
	protoOutDir = filepath.Join(protoOutDir, "type")

	// github.com/foo/bar
	importPath, err := fs.GetImportPath(ctx.RootDir.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// path => baz/qux/quux
	path = strings.Replace(path, "-", "_", -1)

	// quux
	name := filepath.Base(path)

	names := inflect(name)

	// Quux
	serviceName := names.singularCamel

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

	protoPackage := ctx.Config.Package
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

	params := new(params)
	params.Proto.Package = strings.Replace(protoPackage, "-", "_", -1)
	params.PbGo.PackageName = pbgoPackageName
	params.PbGo.PackagePath = filepath.ToSlash(filepath.Join(importPath, pbgoPackagePath))
	params.Name = serviceName
	params.Path = path
	return params, nil
}

type inflectableString struct {
	pluralCamel        string
	pluralCamelLower   string
	pluralSnake        string
	singularCamel      string
	singularCamelLower string
	singularSnake      string
}

func inflect(name string) inflectableString {
	infl := inflectableString{
		pluralCamel:   inflection.Plural(snaker.SnakeToCamel(name)),
		singularCamel: inflection.Singular(snaker.SnakeToCamel(name)),
	}
	infl.pluralCamelLower = strings.ToLower(string(infl.pluralCamel[0])) + infl.pluralCamel[1:]
	infl.pluralSnake = snaker.CamelToSnake(infl.pluralCamel)
	infl.singularCamelLower = strings.ToLower(string(infl.singularCamel[0])) + infl.singularCamel[1:]
	infl.singularSnake = snaker.CamelToSnake(infl.singularCamel)
	return infl
}
