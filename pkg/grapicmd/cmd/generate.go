package cmd

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/serenize/snaker"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate/template"
	"github.com/izumin5210/grapi/pkg/grapicmd/ui"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

func newGenerateCommand(cfg grapicmd.Config, ui ui.UI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate GENERATOR",
		Short:   "Generate new code",
		Aliases: []string{"g", "gen"},
	}

	cmd.AddCommand(newGenerateServiceCommand(cfg, ui))
	cmd.AddCommand(newGenerateCommandCommand(cfg, ui))

	return cmd
}

func newGenerateServiceCommand(cfg grapicmd.Config, ui ui.UI) *cobra.Command {
	return &cobra.Command{
		Use:           "service NAME",
		Short:         "Generate a new service",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			rootDir, ok := fs.LookupRoot(cfg.Fs(), cfg.CurrentDir())
			if !ok {
				return errors.New("geneate command should execut inside a grapi applicaiton directory")
			}

			// github.com/foo/bar
			importPath, err := fs.GetImportPath(rootDir)
			if err != nil {
				return errors.WithStack(err)
			}
			// baz/qux/quux
			path := args[0]

			// quux
			name := filepath.Base(path)
			// Quux
			serviceName := snaker.SnakeToCamel(name)

			// baz/qux
			packagePath := filepath.Dir(path)
			// qux
			packageName := filepath.Base(packagePath)

			// api/baz/qux
			pbgoPackagePath := filepath.Join("api", packagePath)
			// qux_pb
			pbgoPackageName := filepath.Base(pbgoPackagePath) + "_pb"

			if packagePath == "." {
				packagePath = "server"
				packageName = packagePath
				pbgoPackagePath = "api_pb"
				pbgoPackageName = pbgoPackagePath
			}

			protoPackageChunks := []string{}
			for _, pkg := range strings.Split(filepath.Join(importPath, "api", filepath.Dir(path)), "/") {
				chunks := strings.Split(pkg, ".")
				for i := len(chunks) - 1; i >= 0; i-- {
					protoPackageChunks = append(protoPackageChunks, chunks[i])
				}
			}
			// com.github.foo.bar.baz.qux
			protoPackage := strings.Join(protoPackageChunks, ".")

			data := map[string]interface{}{
				"importPath":      importPath,
				"path":            path,
				"name":            name,
				"serviceName":     serviceName,
				"packagePath":     packagePath,
				"packageName":     packageName,
				"pbgoPackagePath": pbgoPackagePath,
				"pbgoPackageName": pbgoPackageName,
				"protoPackage":    protoPackage,
			}
			return generate.NewGenerator(cfg.Fs(), ui, rootDir).Run(template.Service, data)
		},
	}
}

func newGenerateCommandCommand(cfg grapicmd.Config, ui ui.UI) *cobra.Command {
	return &cobra.Command{
		Use:   "command NAME",
		Short: "Generate a new command",
		RunE: func(cmd *cobra.Command, args []string) error {
			rootDir, ok := fs.LookupRoot(cfg.Fs(), cfg.CurrentDir())
			if !ok {
				return errors.New("geneate command should execut inside a grapi applicaiton directory")
			}

			data := map[string]string{
				"name": args[0],
			}
			return generate.NewGenerator(cfg.Fs(), ui, rootDir).Run(template.Command, data)
		},
	}
}
