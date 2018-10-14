package grapicmd

import (
	"io"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"k8s.io/utils/exec"

	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
	"github.com/izumin5210/grapi/pkg/protoc"
)

// Ctx contains the runtime context of grpai.
type Ctx struct {
	FS        afero.Fs
	Viper     *viper.Viper
	Execer    exec.Interface
	InReader  io.Reader
	OutWriter io.Writer
	ErrWriter io.Writer

	CurrentDir string
	RootDir    string
	BinDir     string
	InsideApp  bool

	AppName   string
	Version   string
	Revision  string
	BuildDate string
	Prebuilt  bool

	Config       Config
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

	if c.RootDir == "" {
		c.RootDir, c.InsideApp = fs.LookupRoot(c.FS, c.CurrentDir)
	}

	if c.InsideApp && c.BinDir == "" {
		c.BinDir = filepath.Join(c.RootDir, "bin")
	}
}

// Load reads configurations from the config file.
func (c *Ctx) Load(cfgFile string) error {
	if !c.InsideApp {
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
