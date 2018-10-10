package di

import (
	"sync"

	"github.com/izumin5210/gex"

	"github.com/izumin5210/grapi/pkg/clui"
	"github.com/izumin5210/grapi/pkg/command"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/script"
)

// AppComponent is a dependency container.
type AppComponent interface {
	Config() grapicmd.Config

	UI() clui.UI
	CommandExecutor() command.Executor
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

	ui     clui.UI
	uiOnce sync.Once

	commandExecutor     command.Executor
	commandExecutorOnce sync.Once

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

func (c *appComponentImpl) UI() clui.UI {
	c.uiOnce.Do(func() {
		cfg := c.Config()
		c.ui = clui.New(cfg.OutWriter(), cfg.InReader())
	})
	return c.ui
}

func (c *appComponentImpl) CommandExecutor() command.Executor {
	c.commandExecutorOnce.Do(func() {
		cfg := c.Config()
		c.commandExecutor = command.NewExecutor(cfg.OutWriter(), cfg.ErrWriter(), cfg.InReader())
	})
	return c.commandExecutor
}

func (c *appComponentImpl) ScriptLoader() module.ScriptLoader {
	c.scriptLoaderOnce.Do(func() {
		cfg := c.Config()
		c.scriptLoader = script.NewLoader(cfg.Fs(), c.CommandExecutor(), cfg.RootDir())
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
			WorkingDir: cfg.RootDir(),
			// TODO: set verbose flag
			// TODO: set logger
		}
	})
	return c.gexConfig
}
