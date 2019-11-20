package cmd

import (
	"context"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/izumin5210/clig/pkg/clib"
	"github.com/izumin5210/execx"
	"github.com/izumin5210/gex/pkg/tool"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

func newGenerateCommands(ctx *grapicmd.Ctx) (cmds []*cobra.Command) {
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

	var (
		execs []string
		wg    sync.WaitGroup
	)

	wg.Add(2)
	cmdNames := make(map[string]struct{}, 100)

	go func() {
		defer wg.Done()
		execs = fs.ListExecutableWithPrefix(ctx.FS, "grapi-gen-")
	}()

	go func() {
		defer wg.Done()

		toolRepo, err := di.NewToolRepository(ctx)
		if err != nil {
			zap.L().Debug("failed to initialize tools repository", zap.Error(err))
			return
		}

		tools, err := toolRepo.List(context.TODO())
		if err != nil {
			zap.L().Debug("failed to retrieve tools", zap.Error(err))
			return
		}

		for _, t := range tools {
			if !strings.HasPrefix(t.Name(), "grapi-gen-") {
				continue
			}
			if _, ok := cmdNames[t.Name()]; ok {
				continue
			}
			cmdNames[t.Name()] = struct{}{}
			gCmd.AddCommand(newGenerateCommandByTool(toolRepo, t, "generate"))
			dCmd.AddCommand(newGenerateCommandByTool(toolRepo, t, "destroy"))
		}
	}()

	wg.Wait()

	for _, exec := range execs {
		if _, ok := cmdNames[exec]; ok {
			continue
		}
		cmdNames[exec] = struct{}{}
		gCmd.AddCommand(newGenerateCommandByExec(ctx.IO, ctx.Exec, exec, "generate"))
		dCmd.AddCommand(newGenerateCommandByExec(ctx.IO, ctx.Exec, exec, "destroy"))
	}

	cmds = append(cmds, gCmd, dCmd)

	return
}

func newGenerateCommandByTool(repo tool.Repository, t tool.Tool, subCmd string) *cobra.Command {
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

func newGenerateCommandByExec(io *clib.IO, exec *execx.Executor, path, subCmd string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  strings.TrimPrefix(path, "grapi-gen-"),
		Args: cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			cmd := exec.CommandContext(context.TODO(), path, append([]string{subCmd}, args...)...)
			cmd.Stdout = io.Out
			cmd.Stderr = io.Err
			cmd.Stdin = io.In
			return cmd.Run()
		},
	}
	cmd.SetUsageFunc(func(*cobra.Command) error {
		cmd := exec.CommandContext(context.TODO(), path, subCmd, "--help")
		cmd.Stdout = io.Out
		cmd.Stderr = io.Err
		cmd.Stdin = io.In
		return cmd.Run()
	})
	return cmd
}
