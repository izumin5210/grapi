package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/usecase"
)

type protocCmd *cobra.Command

func provideProtocCommand(cfg *grapicmd.Config, u usecase.ExecuteProtocUsecase) protocCmd {
	return &cobra.Command{
		Use:           "protoc",
		Short:         "Run protoc",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cfg.InsideApp {
				return errors.New("protoc command should be execute inside a grapi application directory")
			}
			return errors.WithStack(u.Perform())
		},
	}
}
