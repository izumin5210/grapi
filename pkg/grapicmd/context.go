package grapicmd

import (
	"os"
	"path/filepath"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/izumin5210/clig/pkg/clib"
	"github.com/izumin5210/execx"
	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/protoc"
)

// Ctx contains the runtime context of grpai.
type Ctx struct {
	FS    afero.Fs
	Viper *viper.Viper
	Exec  *execx.Executor
	IO    *clib.IO

	RootDir   cli.RootDir
	insideApp bool

	Config       Config
	Build        clib.Build
	ProtocConfig protoc.Config
}

// Config stores general setting params and provides accessors for them.
type Config struct {
	Package string
	Grapi   struct {
		ServerDir string
	}
}

// Init initializes the runtime context.
func (c *Ctx) Init() error {
	if c.RootDir.String() == "" {
		dir, _ := os.Getwd()
		c.RootDir = cli.RootDir{clib.Path(dir)}
	}

	if c.IO == nil {
		c.IO = clib.Stdio()
	}

	if c.FS == nil {
		c.FS = afero.NewOsFs()
	}

	if c.Viper == nil {
		c.Viper = viper.New()
	}

	c.Viper.SetFs(c.FS)

	if c.Exec == nil {
		c.Exec = execx.New()
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
		zap.L().Info("failed to find config file", zap.Error(err))
		return nil
	}

	c.insideApp = true
	c.RootDir = cli.RootDir{clib.Path(filepath.Dir(c.Viper.ConfigFileUsed()))}

	err = c.Viper.Unmarshal(&c.Config)
	if err != nil {
		zap.L().Warn("failed to parse config", zap.Error(err))
		return errors.WithStack(err)
	}

	err = c.Viper.UnmarshalKey("protoc", &c.ProtocConfig)
	if err != nil {
		zap.L().Warn("failed to parse protoc config", zap.Error(err))
		return errors.WithStack(err)
	}

	return nil
}

// IsInsideApp returns true if the current working directory is inside a grapi project.
func (c *Ctx) IsInsideApp() bool {
	return c.insideApp
}

// CtxSet is a provider set that includes modules contained in Ctx.
var CtxSet = wire.NewSet(
	ProvideFS,
	ProvideViper,
	ProvideExec,
	ProvideIO,
	ProvideRootDir,
	ProvideConfig,
	ProvideBuildConfig,
	ProvideProtocConfig,
)

func ProvideFS(c *Ctx) afero.Fs                 { return c.FS }
func ProvideViper(c *Ctx) *viper.Viper          { return c.Viper }
func ProvideExec(c *Ctx) *execx.Executor        { return c.Exec }
func ProvideIO(c *Ctx) *clib.IO                 { return c.IO }
func ProvideRootDir(c *Ctx) cli.RootDir         { return c.RootDir }
func ProvideConfig(c *Ctx) *Config              { return &c.Config }
func ProvideBuildConfig(c *Ctx) *clib.Build     { return &c.Build }
func ProvideProtocConfig(c *Ctx) *protoc.Config { return &c.ProtocConfig }
