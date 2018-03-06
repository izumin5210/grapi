package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/command"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
	"github.com/izumin5210/grapi/pkg/grapicmd/ui"
)

func newProtocCommand(cfg grapicmd.Config, ui ui.UI) *cobra.Command {
	return &cobra.Command{
		Use:           "protoc",
		Short:         "Run protoc",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.IsInsideApp() {
				return errors.New("protoc command should be execute inside a grapi applicaiton directory")
			}
			executor := command.NewExecutor(cfg.RootDir(), cfg.OutWriter(), cfg.ErrWriter(), cfg.InReader())
			u := usecase.NewExecuteProtocUsecase(cfg.ProtocConfig(), cfg.Fs(), ui, executor, cfg.RootDir())
			return errors.WithStack(u.Perform())
		},
	}
}
