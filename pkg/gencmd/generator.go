package gencmd

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/izumin5210/clicontrib/pkg/clog"
	"github.com/pkg/errors"
	"github.com/shurcooL/httpfs/vfsutil"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/cli"
)

type Generator interface {
	Generate(params interface{}) error
	Destroy(params interface{}) error
}

func NewGenerator(
	fs afero.Fs,
	ui cli.UI,
	rootDir cli.RootDir,
	templateFS http.FileSystem,
	shouldRunFunc ShouldRunFunc,
) Generator {
	return &generatorImpl{
		fs:            fs,
		ui:            ui,
		rootDir:       rootDir,
		templateFS:    templateFS,
		shouldRunFunc: shouldRunFunc,
	}
}

type generatorImpl struct {
	fs afero.Fs
	ui cli.UI

	rootDir cli.RootDir

	templateFS    http.FileSystem
	shouldRunFunc ShouldRunFunc
}

func (g *generatorImpl) Generate(params interface{}) error {
	entries, err := g.listEntries(params)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, e := range entries {
		if ok, err := g.shouldRun(e); err != nil {
			return errors.WithStack(err)
		} else if !ok {
			continue
		}

		path := g.rootDir.Join(e.Path)

		err := g.fs.MkdirAll(path, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed to create directory")
		}

		err = afero.WriteFile(g.fs, path, []byte(e.Body), 0644)
		if err != nil {
			return errors.Wrapf(err, "failed to write %s", e.Path)
		}
		// TODO: print "Created"
	}

	return nil
}

func (g *generatorImpl) Destroy(params interface{}) error {
	tmplPaths, err := g.listPathTemplates()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, tmplPath := range tmplPaths {
		path, err := TemplateString(strings.TrimSuffix(tmplPath, ".tmpl")).Compile(params)
		if err != nil {
			return errors.Wrapf(err, "failed to parse path: %s", tmplPath)
		}

		paths, err := afero.Glob(g.fs, path)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, path := range paths {
			err = g.fs.Remove(path)
			if err != nil {
				return errors.WithStack(err)
			}
			// TODO: print "Delete"
		}
	}

	return nil
}

func (g *generatorImpl) listEntries(params interface{}) ([]*Entry, error) {
	tmplPaths, err := g.listPathTemplates()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	entries := make([]*Entry, 0, len(tmplPaths))

	for _, tmplPath := range tmplPaths {
		path, err := TemplateString(strings.TrimSuffix(tmplPath, ".tmpl")).Compile(params)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse path: %s", tmplPath)
		}

		data, err := vfsutil.ReadFile(g.templateFS, tmplPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read template: %s", tmplPath)
		}

		body, err := TemplateString(string(data)).Compile(params)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to compile temlpate: %s, %v", tmplPath, params)
		}

		entries = append(entries, &Entry{File: File{Path: filepath.Clean(path), Body: body}, Template: File{Path: tmplPath, Body: string(data)}})
	}
	return entries, nil
}

func (g *generatorImpl) listPathTemplates() (tmplPaths []string, err error) {
	err = vfsutil.Walk(g.templateFS, "/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}

		if info.IsDir() {
			return nil
		}

		tmplPaths = append(tmplPaths, path)

		return nil
	})

	err = errors.WithStack(err)

	return
}

func (g *generatorImpl) shouldRun(e *Entry) (bool, error) {
	if g.shouldRunFunc != nil && !g.shouldRunFunc(e) {
		// TODO: print "Skipped"
		return false, nil
	}

	absPath := g.rootDir.Join(e.Path)

	if ok, err := afero.Exists(g.fs, absPath); err != nil {
		return false, errors.WithStack(err)
	} else if !ok {
		return true, nil
	}

	existed, err := afero.ReadFile(g.fs, absPath)
	if err != nil {
		return false, errors.WithStack(err)
	}

	if string(existed) == e.Body {
		// TODO: print  "Identical"
		return false, nil
	}

	g.ui.ItemFailure(e.Path[1:] + " is conflicted.")
	if ok, err := g.ui.Confirm("Overwite it?"); err != nil {
		clog.Error("failed to confirm to apply", "error", err)
		return false, errors.WithStack(err)
	} else if ok {
		return true, nil
	}

	return false, nil
}
