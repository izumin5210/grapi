package grapicmd

import (
	"io"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// Config stores general setting params and provides accessors for them.
type Config interface {
	Init(cfgFile string)
	Fs() afero.Fs
	AppName() string
	Version() string
	Revision() string
	InReader() io.Reader
	OutWriter() io.Writer
	ErrWriter() io.Writer
}

// NewConfig creates new Config object.
func NewConfig(
	appName, version, revision string,
	in io.Reader,
	out, err io.Writer,
) Config {
	return &config{
		v:        viper.New(),
		fs:       afero.NewOsFs(),
		appName:  appName,
		version:  version,
		revision: revision,
		in:       in,
		out:      out,
		err:      err,
	}
}

type config struct {
	cfgFile                    string
	v                          *viper.Viper
	fs                         afero.Fs
	appName, version, revision string
	in                         io.Reader
	out, err                   io.Writer
	readConfigErr              error
}

func (c *config) Init(cfgFile string) {
	c.v.SetConfigFile(c.cfgFile)
	c.readConfigErr = c.v.ReadInConfig()
}

func (c *config) Fs() afero.Fs {
	return c.fs
}

func (c *config) AppName() string {
	return c.appName
}

func (c *config) Version() string {
	return c.version
}

func (c *config) Revision() string {
	return c.revision
}

func (c *config) InReader() io.Reader {
	return c.in
}

func (c *config) OutWriter() io.Writer {
	return c.out
}

func (c *config) ErrWriter() io.Writer {
	return c.err
}
