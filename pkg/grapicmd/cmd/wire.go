//+build wireinject

package cmd

import (
	"github.com/google/go-cloud/wire"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/di"
)

// NewGrapiCommand creates a new command object.
func NewGrapiCommand(cfg *grapicmd.Config) *cobra.Command {
	wire.Build(
		provideGrapiCommand,
		provideInitCommand,
		provideGenerateCommand,
		provideGenerateServiceCommand,
		provideGenerateScaffoldServiceCommand,
		provideGenerateCommandCommand,
		provideDestroyCommand,
		provideDestroyServiceCommand,
		provideDestroyCommandCommand,
		provideProtocCommand,
		provideBuildCommand,
		provideVersionCommand,
		provideUserDefinedCommands,
		di.Set,
	)
	return nil
}
