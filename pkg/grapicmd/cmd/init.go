package cmd

import (
	"path/filepath"

	"github.com/izumin5210/clicontrib/pkg/clog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
)

var (
	tmplPaths []string
)

func newInitCommand(ctx *grapicmd.Ctx) *cobra.Command {
	var (
		headUsed bool
		pkgName  string
	)

	cmd := &cobra.Command{
		Use:           "init [name]",
		Short:         "Initialize a grapi application",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := parseInitArgs(ctx, args)
			if err != nil {
				return errors.WithStack(err)
			}
			clog.Debug("parseInitArgs", "root", root)

			return errors.WithStack(di.NewInitializeProjectUsecase(ctx).Perform(root, pkgName, headUsed))
		},
	}

	cmd.PersistentFlags().BoolVar(&headUsed, "HEAD", false, "Use HEAD grapi")
	cmd.PersistentFlags().StringVarP(&pkgName, "package", "p", "", `Package name of the application(default: "<parent_package_or_username>.<app_name>")`)

	return cmd
}

func parseInitArgs(ctx *grapicmd.Ctx, args []string) (root string, err error) {
	if argCnt := len(args); argCnt != 1 {
		err = errors.Errorf("invalid argument count: want 1, got %d", argCnt)
		return
	}

	arg := args[0]
	root = ctx.RootDir.String()

	if arg == "." {
		return
	}
	root = arg
	if !filepath.IsAbs(arg) {
		root = ctx.RootDir.Join(arg)
	}
	return
}
