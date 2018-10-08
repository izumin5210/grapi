package di

import (
	"sync"

	"github.com/izumin5210/gex"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/command"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/script"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/ui"
)

// AppComponent is a dependency container.
type AppComponent interface {
	Config() grapicmd.Config

	UI() module.UI
	CommandFactory() module.CommandFactory
	ScriptLoader() module.ScriptLoader
	Generator() module.Generator

	GexConfig() *gex.Config
}

// NewAppComponent creates a new AppComonent instance.
func NewAppComponent(cfg grapicmd.Config) AppComponent {
	return &appComponentImpl{
		config: cfg,
	}
}

type appComponentImpl struct {
	config grapicmd.Config

	ui     module.UI
	uiOnce sync.Once

	commandFactory     module.CommandFactory
	commandFactoryOnce sync.Once

	scriptLoader     module.ScriptLoader
	scriptLoaderOnce sync.Once

	generator     module.Generator
	generatorOnce sync.Once

	gexConfig     *gex.Config
	gexConfigOnce sync.Once
}

func (c *appComponentImpl) Config() grapicmd.Config {
	return c.config
}

func (c *appComponentImpl) UI() module.UI {
	c.uiOnce.Do(func() {
		cfg := c.Config()
		c.ui = ui.New(cfg.OutWriter(), cfg.InReader())
	})
	return c.ui
}

func (c *appComponentImpl) CommandFactory() module.CommandFactory {
	c.commandFactoryOnce.Do(func() {
		cfg := c.Config()
		c.commandFactory = command.NewFactory(cfg.OutWriter(), cfg.ErrWriter(), cfg.InReader())
	})
	return c.commandFactory
}

func (c *appComponentImpl) ScriptLoader() module.ScriptLoader {
	c.scriptLoaderOnce.Do(func() {
		cfg := c.Config()
		c.scriptLoader = script.NewLoader(cfg.Fs(), c.CommandFactory(), cfg.RootDir())
	})
	return c.scriptLoader
}

func (c *appComponentImpl) Generator() module.Generator {
	c.generatorOnce.Do(func() {
		cfg := c.Config()
		c.generator = generator.New(
			cfg.Fs(),
			c.UI(),
			cfg.RootDir(),
			cfg.ProtocConfig().ProtosDir,
			cfg.ProtocConfig().OutDir,
			cfg.ServerDir(),
			cfg.Package(),
			cfg.Version(),
		)
	})
	return c.generator
}

func (c *appComponentImpl) GexConfig() *gex.Config {
	c.gexConfigOnce.Do(func() {
		cfg := c.Config()
		c.gexConfig = &gex.Config{
			OutWriter:  cfg.OutWriter(),
			ErrWriter:  cfg.ErrWriter(),
			InReader:   cfg.InReader(),
			FS:         cfg.Fs(),
			WorkingDir: cfg.CurrentDir(),
			// TODO: set verbose flag
			// TODO: set logger
		}
	})
	return c.gexConfig
}
