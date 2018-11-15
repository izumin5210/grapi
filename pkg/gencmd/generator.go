package gencmd

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/shurcooL/httpfs/vfsutil"
	"github.com/spf13/afero"
	"go.uber.org/zap"

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
			g.ui.ItemFailure(e.Path[1:])
			return errors.WithStack(err)
		} else if !ok {
			g.ui.ItemSkipped(e.Path[1:])
			continue
		}

		err := g.writeFile(e)
		if err != nil {
			g.ui.ItemFailure(e.Path[1:])
			return errors.WithStack(err)
		}
		g.ui.ItemSuccess(e.Path[1:])
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

		absPath := g.rootDir.Join(path)
		if ok, err := afero.Exists(g.fs, absPath); err != nil {
			g.ui.ItemFailure(path)
			return errors.WithStack(err)
		} else if !ok {
			g.ui.ItemSkipped(path)
			continue
		}

		err = g.fs.Remove(absPath)
		if err != nil {
			g.ui.ItemFailure(path)
			return errors.WithStack(err)
		}
		g.ui.ItemSuccess(path)
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
		return false, nil
	}

	g.ui.ItemFailure(e.Path[1:] + " is conflicted.")
	if ok, err := g.ui.Confirm("Overwite it?"); err != nil {
		zap.L().Error("failed to confirm to apply", zap.Error(err))
		return false, errors.WithStack(err)
	} else if ok {
		return true, nil
	}

	return false, nil
}

func (g *generatorImpl) writeFile(e *Entry) error {
	path := g.rootDir.Join(e.Path)
	dir := filepath.Dir(path)

	if ok, err := afero.DirExists(g.fs, dir); err != nil {
		return errors.WithStack(err)
	} else if !ok {
		err := g.fs.MkdirAll(dir, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed to create directory")
		}
	}

	err := afero.WriteFile(g.fs, path, []byte(e.Body), 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write %s", e.Path)
	}

	return nil
}
