package script

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/izumin5210/grapi/pkg/excmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/spf13/afero"
)

type testContext struct {
	fs            afero.Fs
	executor      *excmd.FakeExecutor
	loader        module.ScriptLoader
	binName       string
	rootDir       string
	cmdDir        string
	srcsByBinName map[string][]string
}

func createTestContext(t *testing.T) *testContext {
	fs := afero.NewMemMapFs()
	rootDir := "/home/app"
	binName := "bar"
	ctx := &testContext{
		fs:            fs,
		binName:       binName,
		rootDir:       rootDir,
		cmdDir:        filepath.Join(rootDir, "cmd"),
		executor:      &excmd.FakeExecutor{},
		srcsByBinName: map[string][]string{},
	}

	srcsByBinName := map[string][]string{
		binName: {"bar.go", "foo.go", "main.go"},
	}

	for binName, srcs := range srcsByBinName {
		dir := filepath.Join(ctx.cmdDir, binName)
		err := fs.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		srcPaths := make([]string, 0, len(srcs))
		for _, src := range srcs {
			path := filepath.Join(dir, src)
			err = afero.WriteFile(fs, path, []byte("package main"), 0644)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			srcPaths = append(srcPaths, path)
		}
		ctx.srcsByBinName[binName] = srcPaths
	}

	ctx.loader = NewLoader(fs, ctx.executor, rootDir)

	return ctx
}

func Test_Script(t *testing.T) {
	ctx := createTestContext(t)

	err := ctx.loader.Load(ctx.cmdDir)
	if err != nil {
		t.Errorf("loader.Load() returned an error %v", err)
	}

	s, ok := ctx.loader.Get(ctx.binName)
	if got, want := ok, true; got != want {
		t.Errorf("loader.Get() returned %t, want %t", got, want)
	}

	if got, want := s.Name(), ctx.binName; got != want {
		t.Errorf("script.Name() returned %v, want %v", got, want)
	}

	if err != nil {
		t.Errorf("script.Build() returned an error %v", err)
	}

	err = s.Build("-v")
	binPath := filepath.Join(ctx.rootDir, "bin", ctx.binName)
	cmd := ctx.executor.Commands[0]

	srcs := ctx.srcsByBinName[ctx.binName]
	if got, want := cmd.Args, append([]string{"build", "-o=" + binPath, "-v"}, srcs...); !reflect.DeepEqual(got, want) {
		t.Errorf("Build() cmduted %v, want %v", got, want)
	}

	if got, want := cmd.Dir, "/home/app"; got != want {
		t.Errorf("Build() cmduted a command in %v, want %v", got, want)
	}

	if err != nil {
		t.Errorf("Build() returned an error %v", err)
	}

	err = s.Run("-v")
	cmd = ctx.executor.Commands[1]

	if got, want := cmd.Name, binPath; got != want {
		t.Errorf("Run() cmduted %v, want %v", got, want)
	}

	if got, want := cmd.Dir, "/home/app"; got != want {
		t.Errorf("Run() cmduted a command in %v, want %v", got, want)
	}

	if err != nil {
		t.Errorf("Run() returned an error %v", err)
	}
}
