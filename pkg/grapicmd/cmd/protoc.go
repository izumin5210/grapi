package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
)

func newProtocCommand(ctx *grapicmd.Ctx) *cobra.Command {
	return &cobra.Command{
		Use:           "protoc",
		Short:         "Run protoc",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ctx.InsideApp {
				return errors.New("protoc command should be execute inside a grapi application directory")
			}
			return errors.WithStack(di.NewExecuteProtocUsecase(ctx).Perform())
		},
	}
}
