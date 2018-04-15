package script

import (
	"path/filepath"
	"runtime"
	"sort"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

// NewLoader creates a new ScriptLoader instance.
func NewLoader(fs afero.Fs, commandFactory module.CommandFactory, rootDir string) module.ScriptLoader {
	return &scriptLoader{
		fs:             fs,
		commandFactory: commandFactory,
		rootDir:        rootDir,
		binDir:         filepath.Join(rootDir, "bin"),
		scripts:        make(map[string]module.Script),
	}
}

type scriptLoader struct {
	fs             afero.Fs
	commandFactory module.CommandFactory
	rootDir        string
	binDir         string
	scripts        map[string]module.Script
	names          []string
}

func (f *scriptLoader) Load(dir string) error {
	srcsByDir, err := fs.FindMainPackagesAndSources(f.fs, dir)
	if err != nil {
		return errors.Wrap(err, "failed to find commands")
	}
	for dir, srcs := range srcsByDir {
		srcPaths := make([]string, 0, len(srcs))
		for _, name := range srcs {
			srcPaths = append(srcPaths, filepath.Join(dir, name))
		}
		name := filepath.Base(dir)
		ext := ""
		if runtime.GOOS == "windows" {
			ext = ".exe"
		}
		f.scripts[name] = &script{
			fs:             f.fs,
			commandFactory: f.commandFactory,
			srcPaths:       srcPaths,
			name:           name,
			binPath:        filepath.Join(f.binDir, name+ext),
			rootDir:        f.rootDir,
		}
		f.names = append(f.names, name)
	}
	sort.Strings(f.names)
	return nil
}

func (f *scriptLoader) Get(name string) (script module.Script, ok bool) {
	script, ok = f.scripts[name]
	return
}

func (f *scriptLoader) Names() []string {
	return f.names
}
