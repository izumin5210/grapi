package gencmd

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Command represents a subcommand of a generator plugin. It will be converted to a *cobra.Command object internally.
type Command struct {
	// Use, Short, Long, Example and Args are pass-through into *cobra.Command object.
	Use     string
	Short   string
	Long    string
	Example string
	Args    cobra.PositionalArgs

	// BuildParams returns parameters to generate/destroy files(required).
	BuildParams func(c *Command, args []string) (interface{}, error)

	// PreRun is executed in *cobra.Command.PreRunE.
	PreRun func(c *Command, args []string) error

	// PostRun is executed in *cobra.Command.PostRunE.
	PostRun func(c *Command, args []string) error

	// ShouldRun is executed for each generated files. When it returns false, the file will be skipped.
	ShouldRun ShouldRunFunc

	// ShouldInsideApp will disable the command when a current working directory is not inside of a grapi project.
	ShouldInsideApp bool

	// TemplateFS contains file templates(required).
	TemplateFS http.FileSystem

	flags *pflag.FlagSet
	ctx   *Ctx
}

// Flags returns a FlagSet that applies to this commmand.
func (c *Command) Flags() *pflag.FlagSet {
	if c.flags == nil {
		c.flags = new(pflag.FlagSet)
	}
	return c.flags
}

// Ctx returns the context object.
func (c *Command) Ctx() *Ctx {
	return c.ctx
}

func (c *Command) newCobraCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:     c.Use,
		Short:   c.Short,
		Long:    c.Long,
		Example: c.Example,
		Args:    c.Args,
		PreRunE: func(_ *cobra.Command, args []string) error {
			if c.ShouldInsideApp && !c.Ctx().IsInsideApp() {
				return errors.New("should execute inside grapi project")
			}
			if c.PreRun != nil {
				err := c.PreRun(c, args)
				if err != nil {
					return errors.WithStack(err)
				}
			}
			return nil
		},
	}
	if c.PostRun != nil {
		cc.PostRunE = func(_ *cobra.Command, args []string) error { return c.PostRun(c, args) }
	}
	cc.PersistentFlags().AddFlagSet(c.Flags())
	return cc
}

// File represents a file content.
type File struct {
	Path string
	Body string
}

// Entry represents a file that will be generated.
type Entry struct {
	File
	Template File
}

type ShouldRunFunc func(e *Entry) bool
