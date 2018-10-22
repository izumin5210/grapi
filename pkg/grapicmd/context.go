package grapicmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"k8s.io/utils/exec"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/protoc"
)

// Ctx contains the runtime context of grpai.
type Ctx struct {
	FS     afero.Fs
	Viper  *viper.Viper
	Execer exec.Interface
	IO     *cli.IO

	RootDir   cli.RootDir
	insideApp bool

	Config       Config
	Build        BuildConfig
	ProtocConfig protoc.Config
}

// Config stores general setting params and provides accessors for them.
type Config struct {
	Package string
	Grapi   struct {
		ServerDir string
	}
}

type BuildConfig struct {
	AppName   string
	Version   string
	Revision  string
	BuildDate string
	Prebuilt  bool
}

// Init initializes the runtime context.
func (c *Ctx) Init() {
	if c.FS == nil {
		c.FS = afero.NewOsFs()
	}

	if c.Viper == nil {
		c.Viper = viper.New()
		c.Viper.SetFs(c.FS)
	}

	if c.Execer == nil {
		c.Execer = exec.New()
	}

	cwd := c.RootDir
	c.RootDir, c.insideApp = cli.LookupRoot(c.FS, string(cwd))
	if c.RootDir == "" {
		c.RootDir = cwd
	}
}

// Load reads configurations from the config file.
func (c *Ctx) Load(cfgFile string) error {
	if !c.IsInsideApp() {
		return nil
	}

	c.Viper.SetConfigFile(cfgFile)
	err := c.Viper.ReadInConfig()
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.Viper.UnmarshalKey("config", &c.Config)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.Viper.UnmarshalKey("protoc", &c.ProtocConfig)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Ctx) IsInsideApp() bool {
	return c.insideApp
}
