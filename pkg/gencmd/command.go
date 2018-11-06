package gencmd

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CommandType int

const (
	CommandUnknown CommandType = iota
	CommandGenerate
	CommandDestroy
)

type Command struct {
	Use             string
	Short           string
	Long            string
	Example         string
	Args            cobra.PositionalArgs
	BuildParams     func(c *Command, args []string) (interface{}, error)
	PreRun          func(c *Command, args []string) error
	PostRun         func(c *Command, args []string) error
	ShouldRun       ShouldRunFunc
	ShouldInsideApp bool
	TemplateFS      http.FileSystem

	flags *pflag.FlagSet
	ctx   *Ctx
}

func (c *Command) Flags() *pflag.FlagSet {
	if c.flags == nil {
		c.flags = new(pflag.FlagSet)
	}
	return c.flags
}

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

type File struct {
	Path string
	Body string
}

type Entry struct {
	File
	Template File
}

type ShouldRunFunc func(e *Entry) bool
