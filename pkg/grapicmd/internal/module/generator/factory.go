package generator

import (
	assets "github.com/jessevdk/go-assets"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator/template"
)

// NewFactory creates a new module.GeneratorFactory instance.
func NewFactory(fs afero.Fs, ui module.UI) module.GeneratorFactory {
	return &factory{
		fs: fs,
		ui: ui,
	}
}

type factory struct {
	fs afero.Fs
	ui module.UI
}

func (f *factory) Project() module.Generator {
	return f.create(template.Init)
}

func (f *factory) Service() module.Generator {
	return f.create(template.Service)
}

func (f *factory) Command() module.Generator {
	return f.create(template.Command)
}

func (f *factory) create(tmplFs *assets.FileSystem) module.Generator {
	return &generator{
		tmplFs: tmplFs,
		fs:     f.fs,
		ui:     f.ui,
	}
}
