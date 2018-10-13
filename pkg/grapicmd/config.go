package grapicmd

import (
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/izumin5210/grapi/pkg/grapicmd/protoc"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// Config stores general setting params and provides accessors for them.
type Config struct {
	Fs           afero.Fs
	CurrentDir   string
	RootDir      string
	InsideApp    bool
	AppName      string
	Version      string
	Revision     string
	BuildDate    string
	Prebuilt     bool
	InReader     io.Reader
	OutWriter    io.Writer
	ErrWriter    io.Writer
	ServerDir    string
	Package      string
	ProtocConfig *protoc.Config

	viper *viper.Viper
}

// NewConfig creates new Config object.
func NewConfig(
	currentDir string,
	appName, version string,
	revision, buildDate string,
	prebuilt bool,
	in io.Reader,
	out, err io.Writer,
) *Config {
	afs := afero.NewOsFs()
	v := viper.New()
	v.SetFs(afs)
	rootDir, insideApp := fs.LookupRoot(afs, currentDir)
	return &Config{
		Fs:           afs,
		CurrentDir:   currentDir,
		RootDir:      rootDir,
		InsideApp:    insideApp,
		AppName:      appName,
		Version:      version,
		Revision:     revision,
		BuildDate:    buildDate,
		Prebuilt:     prebuilt,
		InReader:     in,
		OutWriter:    out,
		ErrWriter:    err,
		ProtocConfig: &protoc.Config{},
		viper:        v,
	}
}

func (c *Config) Init(cfgFile string) error {
	if !c.InsideApp {
		return nil
	}

	c.viper.SetConfigFile(cfgFile)
	err := c.viper.ReadInConfig()
	if err != nil {
		return errors.WithStack(err)
	}
	c.Package = c.viper.GetString("package")
	c.ServerDir = c.viper.GetString("grapi.server_dir")
	err = c.viper.UnmarshalKey("protoc", c.ProtocConfig)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
