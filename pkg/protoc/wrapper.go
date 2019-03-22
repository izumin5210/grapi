package protoc

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/izumin5210/gex/pkg/tool"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"k8s.io/utils/exec"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// Wrapper can execute protoc commands for current project's proto files.
type Wrapper interface {
	Exec(context.Context) error
}

type wrapperImpl struct {
	cfg      *Config
	fs       afero.Fs
	ui       cli.UI
	execer   exec.Interface
	toolRepo tool.Repository
	rootDir  cli.RootDir
	modCache *sync.Map
}

// NewWrapper creates a new Wrapper instance.
func NewWrapper(cfg *Config, fs afero.Fs, execer exec.Interface, ui cli.UI, toolRepo tool.Repository, rootDir cli.RootDir) Wrapper {
	return &wrapperImpl{
		cfg:      cfg,
		fs:       fs,
		ui:       ui,
		execer:   execer,
		toolRepo: toolRepo,
		rootDir:  rootDir,
		modCache: new(sync.Map),
	}
}

func (e *wrapperImpl) Exec(ctx context.Context) (err error) {
	e.ui.Section("Execute protoc")

	e.ui.Subsection("Install plugins")
	err = errors.WithStack(e.installPlugins(ctx))
	if err != nil {
		return
	}

	e.ui.Subsection("Execute protoc")
	err = errors.WithStack(e.execProtocAll(ctx))

	return
}

func (e *wrapperImpl) installPlugins(ctx context.Context) error {
	return errors.WithStack(e.toolRepo.BuildAll(ctx))
}

func (e *wrapperImpl) execProtocAll(ctx context.Context) error {
	protoFiles, err := e.cfg.ProtoFiles(e.fs, e.rootDir.String())
	if err != nil {
		return errors.WithStack(err)
	}

	var errs []error
	for _, path := range protoFiles {
		err = e.execProtoc(ctx, path)
		relPath, _ := filepath.Rel(e.rootDir.String(), path)
		if err == nil {
			e.ui.ItemSuccess(relPath)
		} else {
			zap.L().Error("failed to execute protoc", zap.Error(err))
			errs = append(errs, err)
			e.ui.ItemFailure(relPath, err)
		}
	}

	if len(errs) > 0 {
		return errors.New("failed to execute protoc")
	}

	return nil
}

func (e *wrapperImpl) execProtoc(ctx context.Context, protoPath string) error {
	outDir, err := e.cfg.OutDirOf(e.rootDir.String(), protoPath)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = fs.CreateDirIfNotExists(e.fs, outDir); err != nil {
		return errors.WithStack(err)
	}

	cmds, err := e.commands(ctx, protoPath)
	if err != nil {
		return errors.WithStack(err)
	}

	path := e.rootDir.BinDir().String() + string(filepath.ListSeparator) + os.Getenv("PATH")
	env := append(os.Environ(), "PATH="+path)

	for _, args := range cmds {
		cmd := e.execer.CommandContext(ctx, args[0], args[1:]...)
		cmd.SetEnv(env)
		cmd.SetDir(e.rootDir.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			return errors.Wrapf(err, "failed to execute command: %v\n%s", args, string(out))
		}
	}

	return nil
}

func (e *wrapperImpl) commands(ctx context.Context, protoPath string) ([][]string, error) {
	cmds := make([][]string, 0, len(e.cfg.Plugins))
	relProtoPath, _ := filepath.Rel(e.rootDir.String(), protoPath)

	outDir, err := e.cfg.OutDirOf(e.rootDir.String(), protoPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	isMod := e.modulesEnabled()

	funcMap := template.FuncMap{
		"module": func(in string) (string, error) {
			if isMod {
				return e.getModulePath(ctx, in)
			}
			return e.rootDir.Join("vendor", in).String(), nil
		},
	}

	for _, p := range e.cfg.Plugins {
		args := []string{}
		args = append(args, "-I", filepath.Dir(relProtoPath))
		for _, dir := range e.cfg.ImportDirs {
			tmpl, err := template.New(dir).Funcs(funcMap).Parse(dir)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			buf := new(bytes.Buffer)
			err = tmpl.Funcs(funcMap).Execute(buf, nil)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			args = append(args, "-I", buf.String())
		}
		args = append(args, p.toProtocArg(outDir))
		args = append(args, relProtoPath)
		cmds = append(cmds, append([]string{"protoc"}, args...))
	}

	return cmds, nil
}

func (e *wrapperImpl) getModulePath(ctx context.Context, pkg string) (string, error) {
	if v, ok := e.modCache.Load(pkg); ok {
		return v.(string), nil
	}
	buf := new(bytes.Buffer)
	cmd := e.execer.CommandContext(ctx, "go", "list", "-f", "{{.Dir}}", "-m", pkg)
	cmd.SetEnv(append(os.Environ(), "GO111MODULE", "on"))
	cmd.SetDir(e.rootDir.String())
	cmd.SetStdout(buf)
	err := cmd.Run()
	if err != nil {
		return "", errors.WithStack(err)
	}
	out := strings.TrimSpace(buf.String())
	e.modCache.Store(pkg, out)

	return out, nil
}

func (e *wrapperImpl) modulesEnabled() bool {
	for _, f := range []string{"go.mod", "go.sum"} {
		if ok, err := afero.Exists(e.fs, e.rootDir.Join(f).String()); err != nil || !ok {
			return false
		}
	}
	return true
}
