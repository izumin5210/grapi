package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/izumin5210/clicontrib/pkg/clog"
	"github.com/izumin5210/gex/pkg/tool"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
)

func newGenerateCommands(ctx *grapicmd.Ctx) (cmds []*cobra.Command) {
	toolRepo, err := di.NewToolRepository(ctx)
	if err != nil {
		clog.Debug("failed to initialize tools repository", "error", err)
		return
	}

	tools, err := toolRepo.List(context.TODO())
	if err != nil {
		clog.Debug("failed to retrieve tools", "error", err)
		return
	}

	gCmd := &cobra.Command{
		Use:     "generate GENERATOR",
		Short:   "Generate a new code",
		Aliases: []string{"g", "gen"},
	}
	dCmd := &cobra.Command{
		Use:     "destroy GENERATOR",
		Short:   "Destroy an existing new code",
		Aliases: []string{"d"},
	}

	for _, t := range tools {
		if strings.HasPrefix(t.Name(), "grapi-gen-") {
			gCmd.AddCommand(newGenerateCommand(toolRepo, t, "generate"))
			dCmd.AddCommand(newGenerateCommand(toolRepo, t, "destroy"))
		}
	}

	cmds = append(cmds, gCmd, dCmd)

	return
}

func newGenerateCommand(repo tool.Repository, t tool.Tool, subCmd string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  strings.TrimPrefix(t.Name(), "grapi-gen-"),
		Args: cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return repo.Run(context.TODO(), t.Name(), append([]string{subCmd}, args...)...)
		},
	}
	cmd.SetUsageFunc(func(*cobra.Command) error {
		return repo.Run(context.TODO(), t.Name(), subCmd, "--help")
	})
	return cmd
}
