package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd/di"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
)

func newProtocCommand(ac di.AppComponent) *cobra.Command {
	return &cobra.Command{
		Use:           "protoc",
		Short:         "Run protoc",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ac.Config()
			if !cfg.InsideApp {
				return errors.New("protoc command should be execute inside a grapi application directory")
			}
			u := usecase.NewExecuteProtocUsecase(cfg.ProtocConfig, cfg.Fs, ac.UI(), ac.CommandExecutor(), ac.GexConfig(), cfg.RootDir)
			return errors.WithStack(u.Perform())
		},
	}
}
