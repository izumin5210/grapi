package cmd

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/command"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate"
	"github.com/izumin5210/grapi/pkg/grapicmd/generate/template"
	"github.com/izumin5210/grapi/pkg/grapicmd/ui"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

var (
	tmplPaths []string
)

func newInitCommand(cfg grapicmd.Config, ui ui.UI) *cobra.Command {
	var (
		depSkipped bool
	)

	cmd := &cobra.Command{
		Use:           "init [name]",
		Short:         "Initialize a grapi application",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			root := cfg.CurrentDir()

			if argCnt := len(args); argCnt == 1 {
				arg := args[0]
				if arg != "." {
					if filepath.IsAbs(arg) {
						root = arg
					} else {
						root, err = filepath.Abs(arg)
						if err != nil {
							return errors.Wrap(err, "failed to get the target directory")
						}
					}
				}
			} else if argCnt > 1 {
				return errors.Errorf("invalid argument count: want 0 or 1, got %d", argCnt)
			}

			ui.Section("Initialize project")
			err = initProject(cfg.Fs(), ui, root)
			if err != nil {
				return errors.Wrap(err, "failed to initialize project")
			}

			if !depSkipped {
				ui.Subsection("Install dependencies")
				cmdExecutor := command.NewExecutor(root, cfg.OutWriter(), cfg.ErrWriter(), cfg.InReader())
				_, err = cmdExecutor.Exec(
					[]string{"dep", "ensure", "-v"},
					command.WithIOConnected(),
				)
				err = errors.Wrapf(err, "failed to execute `dep ensure`")
			}

			return err
		},
	}

	cmd.PersistentFlags().BoolVarP(&depSkipped, "skip-dep", "D", false, "Don't run `dep ensure`")

	return cmd
}

func initProject(afs afero.Fs, ui ui.UI, rootPath string) error {
	importPath, err := fs.GetImportPath(rootPath)
	if err != nil {
		return errors.WithStack(err)
	}
	data := map[string]string{
		"importPath": importPath,
	}
	return generate.NewGenerator(afs, ui, rootPath).Run(template.Init, data)
}
