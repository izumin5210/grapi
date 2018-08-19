package grapicmd

import (
	"io"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/izumin5210/grapi/pkg/grapicmd/protoc"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// Config stores general setting params and provides accessors for them.
type Config interface {
	Init(cfgFile string)
	Fs() afero.Fs
	CurrentDir() string
	RootDir() string
	IsInsideApp() bool
	AppName() string
	Version() string
	Revision() string
	BuildDate() string
	ReleaseType() string
	InReader() io.Reader
	OutWriter() io.Writer
	ErrWriter() io.Writer
	ServerDir() string
	Package() string
	ProtocConfig() *protoc.Config
}

// NewConfig creates new Config object.
func NewConfig(
	currentDir string,
	appName, version, revision string,
	buildDate, releaseType string,
	in io.Reader,
	out, err io.Writer,
) Config {
	afs := afero.NewOsFs()
	rootDir, insideApp := fs.LookupRoot(afs, currentDir)
	return &config{
		v:           viper.New(),
		fs:          afs,
		currentDir:  currentDir,
		rootDir:     rootDir,
		insideApp:   insideApp,
		appName:     appName,
		version:     version,
		revision:    revision,
		buildDate:   buildDate,
		releaseType: releaseType,
		in:          in,
		out:         out,
		err:         err,
	}
}

type config struct {
	cfgFile                    string
	v                          *viper.Viper
	fs                         afero.Fs
	currentDir, rootDir        string
	insideApp                  bool
	appName, version, revision string
	buildDate, releaseType     string
	in                         io.Reader
	out, err                   io.Writer
	readConfigErr              error
}

func (c *config) Init(cfgFile string) {
	c.cfgFile = cfgFile
	c.v.SetConfigFile(c.cfgFile)
	c.readConfigErr = c.v.ReadInConfig()
}

func (c *config) Fs() afero.Fs {
	return c.fs
}

func (c *config) CurrentDir() string {
	return c.currentDir
}

func (c *config) RootDir() string {
	return c.rootDir
}

func (c *config) IsInsideApp() bool {
	return c.insideApp
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

func (c *config) BuildDate() string {
	return c.buildDate
}

func (c *config) ReleaseType() string {
	return c.releaseType
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

func (c *config) Package() string {
	return c.v.GetString("package")
}

func (c *config) ServerDir() string {
	return c.v.GetString("grapi.server_dir")
}

func (c *config) ProtocConfig() *protoc.Config {
	cfg := &protoc.Config{}
	c.v.UnmarshalKey("protoc", cfg)
	return cfg
}
