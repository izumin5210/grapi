package cmd

import (
	"path/filepath"

	"github.com/izumin5210/clicontrib/clog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
)

var (
	tmplPaths []string
)

func newInitCommand(cfg grapicmd.Config, ui module.UI, generator module.ProjectGenerator, commandFactory module.CommandFactory) *cobra.Command {
	var (
		depSkipped bool
		headUsed   bool
	)

	cmd := &cobra.Command{
		Use:           "init [name]",
		Short:         "Initialize a grapi application",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := parseInitArgs(cfg, args)
			if err != nil {
				return errors.WithStack(err)
			}
			clog.Debug("parseInitArgs", "root", root)

			u := usecase.NewInitializeProjectUsecase(
				ui,
				generator,
				commandFactory,
				cfg.Version(),
			)

			return errors.WithStack(u.Perform(root, depSkipped, headUsed))
		},
	}

	cmd.PersistentFlags().BoolVarP(&depSkipped, "skip-dep", "D", false, "Don't run dep ensure")
	cmd.PersistentFlags().BoolVar(&headUsed, "HEAD", false, "Use HEAD grapi")

	return cmd
}

func parseInitArgs(cfg grapicmd.Config, args []string) (root string, err error) {
	if argCnt := len(args); argCnt != 1 {
		err = errors.Errorf("invalid argument count: want 1, got %d", argCnt)
		return
	}

	arg := args[0]
	root = cfg.CurrentDir()

	if arg == "." {
		return
	}
	root = arg
	if !filepath.IsAbs(arg) {
		root = filepath.Join(cfg.CurrentDir(), arg)
	}
	return
}
