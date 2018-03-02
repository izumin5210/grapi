package generate

import (
	"path/filepath"
	"sort"
	"strings"

	assets "github.com/jessevdk/go-assets"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/grapicmd/ui"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// Generator is an interface to generate files from template and given params.
type Generator interface {
	Run(data interface{}) error
}

// NewGenerator createes a new generator instance with specified filessytem and templates.
func NewGenerator(fs afero.Fs, ui ui.UI, rootPath string, tmplFs *assets.FileSystem) Generator {
	return &generator{
		fs:       fs,
		ui:       ui,
		rootPath: rootPath,
		tmplFs:   tmplFs,
	}
}

type generator struct {
	fs       afero.Fs
	ui       ui.UI
	rootPath string
	tmplFs   *assets.FileSystem
}

func (g *generator) Run(data interface{}) error {
	for _, tmplPath := range g.sortedEntryPaths() {
		entry := g.tmplFs.Files[tmplPath]
		path, err := TemplateString(strings.TrimSuffix(tmplPath, ".tmpl")).Compile(data)
		if err != nil {
			return errors.Wrapf(err, "failed to parse path: %s", path)
		}
		absPath := filepath.Join(g.rootPath, path)
		dirPath := filepath.Dir(absPath)

		// create directory if not exists
		if err := fs.CreateDirIfNotExists(g.fs, dirPath); err != nil {
			return errors.WithStack(err)
		}

		// generate content
		body, err := TemplateString(string(entry.Data)).Compile(data)
		if err != nil {
			return errors.Wrapf(err, "failed to generate %s", path)
		}

		// check existed entries
		st := statusCreate
		if ok, err := afero.Exists(g.fs, absPath); err != nil {
			// TODO: handle an error
			st = statusSkipped
		} else if ok {
			existedBody, err := afero.ReadFile(g.fs, absPath)
			if err != nil {
				// TODO: handle an error
				st = statusSkipped
			}
			if string(existedBody) == body {
				st = statusIdentical
			} else {
				// TODO: ask to overwrite
				st = statusConflicted
			}
		}

		// create
		if st.ShouldCreate() {
			err = afero.WriteFile(g.fs, absPath, []byte(body), 0644)
			if err != nil {
				return errors.Wrapf(err, "failed to write %s", path)
			}
		}

		g.ui.PrintWithStatus(path[1:], st)
	}

	return nil
}

func (g *generator) sortedEntryPaths() []string {
	rootFiles := make([]string, 0, len(g.tmplFs.Files))
	tmplPaths := make([]string, 0, len(g.tmplFs.Files))
	for path, entry := range g.tmplFs.Files {
		if entry.IsDir() {
			continue
		}
		if strings.Count(entry.Path[1:], "/") == 0 {
			rootFiles = append(rootFiles, path)
		} else {
			tmplPaths = append(tmplPaths, path)
		}
	}
	sort.Strings(rootFiles)
	sort.Strings(tmplPaths)
	return append(rootFiles, tmplPaths...)
}
