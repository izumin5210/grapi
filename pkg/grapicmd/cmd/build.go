package cmd

import (
	"path/filepath"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal"
	"github.com/izumin5210/grapi/pkg/grapicmd/ui"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newBuildCommand(cfg grapicmd.Config, ui ui.UI, scriptFactory internal.ScriptFactory) *cobra.Command {
	return &cobra.Command{
		Use:           "build [TARGET]...",
		Short:         "Build commands",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(c *cobra.Command, args []string) (err error) {
			rootDir, ok := fs.LookupRoot(cfg.Fs(), cfg.CurrentDir())
			if !ok {
				return errors.New("protoc command should be execute inside a grapi applicaiton directory")
			}

			nameSet := make(map[string]bool, len(args))
			for _, n := range args {
				nameSet[n] = true
			}
			isAll := len(args) == 0

			paths, err := afero.Glob(cfg.Fs(), filepath.Join(rootDir, "cmd/*/run.go"))
			for _, path := range paths {
				script := scriptFactory.Create(path)
				if isAll || nameSet[script.Name()] {
					ui.Subsection("Building " + script.Name())
					err := script.Build()
					if err != nil {
						return errors.WithStack(err)
					}
				}
			}

			return nil
		},
	}
}
