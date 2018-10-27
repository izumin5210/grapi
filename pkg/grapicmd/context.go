package grapicmd

import (
	"path/filepath"

	"github.com/izumin5210/clicontrib/pkg/clog"
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
func (c *Ctx) Init() error {
	if c.FS == nil {
		c.FS = afero.NewOsFs()
	}

	if c.Viper == nil {
		c.Viper = viper.New()
	}

	c.Viper.SetFs(c.FS)

	if c.Execer == nil {
		c.Execer = exec.New()
	}

	if c.Build.AppName == "" {
		c.Build.AppName = "grapi"
	}

	return errors.WithStack(c.loadConfig())
}

func (c *Ctx) loadConfig() error {
	c.Viper.SetConfigName(c.Build.AppName)
	for dir := c.RootDir.String(); dir != "/"; dir = filepath.Dir(dir) {
		c.Viper.AddConfigPath(dir)
	}

	err := c.Viper.ReadInConfig()
	if err != nil {
		clog.Info("failed to find config file", "error", err)
		return nil
	}

	c.insideApp = true
	c.RootDir = cli.RootDir(filepath.Dir(c.Viper.ConfigFileUsed()))

	err = c.Viper.Unmarshal(&c.Config)
	if err != nil {
		clog.Warn("failed to parse config", "error", err)
		return errors.WithStack(err)
	}

	err = c.Viper.UnmarshalKey("protoc", &c.ProtocConfig)
	if err != nil {
		clog.Warn("failed to parse protoc config", "error", err)
		return errors.WithStack(err)
	}

	return nil
}

func (c *Ctx) IsInsideApp() bool {
	return c.insideApp
}
