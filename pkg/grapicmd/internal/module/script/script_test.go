package script

import (
	"context"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/afero"

	"github.com/izumin5210/clig/pkg/clib"
	"github.com/izumin5210/execx"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
)

type testContext struct {
	fs            afero.Fs
	loader        module.ScriptLoader
	binName       string
	rootDir       string
	cmdDir        string
	srcsByBinName map[string][]string
	cmds          []*exec.Cmd
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
		srcsByBinName: map[string][]string{},
	}
	exec := execx.New(execx.WithFakeProcess(
		func(_ context.Context, c *exec.Cmd) error { ctx.cmds = append(ctx.cmds, c); return nil },
	))

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

	ctx.loader = NewLoader(fs, &clib.IO{}, exec, rootDir)

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

	err = s.Build(context.Background(), "-v")
	binPath := filepath.Join(ctx.rootDir, "bin", ctx.binName)
	cmd := ctx.cmds[0]

	srcs := ctx.srcsByBinName[ctx.binName]
	if got, want := cmd.Args, append([]string{"go", "build", "-o=" + binPath, "-v"}, srcs...); !reflect.DeepEqual(got, want) {
		t.Errorf("Build() executed %v, want %v", got, want)
	}

	if got, want := cmd.Dir, "/home/app"; got != want {
		t.Errorf("Build() executed a command in %v, want %v", got, want)
	}

	if err != nil {
		t.Errorf("Build() returned an error %v", err)
	}

	err = s.Run(context.Background(), "-v")
	cmd = ctx.cmds[1]

	if got, want := cmd.Path, binPath; got != want {
		t.Errorf("Run() executed %v, want %v", got, want)
	}

	if got, want := cmd.Dir, "/home/app"; got != want {
		t.Errorf("Run() executed a command in %v, want %v", got, want)
	}

	if err != nil {
		t.Errorf("Run() returned an error %v", err)
	}
}
