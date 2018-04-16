package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
)

func newProtocCommand(cfg grapicmd.Config, ui module.UI, commandFactory module.CommandFactory) *cobra.Command {
	return &cobra.Command{
		Use:           "protoc",
		Short:         "Run protoc",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.IsInsideApp() {
				return errors.New("protoc command should be execute inside a grapi application directory")
			}
			u := usecase.NewExecuteProtocUsecase(cfg.ProtocConfig(), cfg.Fs(), ui, commandFactory, cfg.RootDir())
			return errors.WithStack(u.Perform())
		},
	}
}
