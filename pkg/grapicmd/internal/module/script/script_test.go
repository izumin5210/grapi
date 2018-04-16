package script

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/spf13/afero"
)

type execution struct {
	nameAndArgs []string
	dir         string
	connected   bool
}

type testContext struct {
	fs            afero.Fs
	executions    []*execution
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
		executions:    []*execution{},
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

	commandFactory := &fakeCommandFactory{
		fakeCreate: func(nameAndArgs []string) module.Command {
			c := &fakeCommand{}
			c.fakeExec = func() ([]byte, error) {
				ctx.executions = append(ctx.executions, &execution{
					nameAndArgs: nameAndArgs,
					dir:         c.dir,
					connected:   c.ioConnected,
				})
				return []byte{}, nil
			}
			return c
		},
	}

	ctx.loader = NewLoader(fs, commandFactory, rootDir)

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
	exec := ctx.executions[0]

	srcs := ctx.srcsByBinName[ctx.binName]
	if got, want := exec.nameAndArgs, append([]string{"go", "build", "-o=" + binPath, "-v"}, srcs...); !reflect.DeepEqual(got, want) {
		t.Errorf("Build() executed %v, want %v", got, want)
	}

	if got, want := exec.dir, "/home/app"; got != want {
		t.Errorf("Build() executed a command in %v, want %v", got, want)
	}

	if err != nil {
		t.Errorf("Build() returned an error %v", err)
	}

	err = s.Run("-v")
	exec = ctx.executions[1]

	if got, want := exec.nameAndArgs, []string{binPath, "-v"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Run() executed %v, want %v", got, want)
	}

	if got, want := exec.dir, "/home/app"; got != want {
		t.Errorf("Run() executed a command in %v, want %v", got, want)
	}
	
	if err != nil {
		t.Errorf("Run() returned an error %v", err)
	}
}

// fake impls

type fakeCommandFactory struct {
	module.CommandFactory
	fakeCreate func([]string) module.Command
}

func (f *fakeCommandFactory) Create(nameAndArgs []string) module.Command {
	return f.fakeCreate(nameAndArgs)
}

type fakeCommand struct {
	module.Command
	dir         string
	ioConnected bool
	fakeExec    func() ([]byte, error)
}

func (c *fakeCommand) SetDir(dir string) module.Command {
	c.dir = dir
	return c
}

func (c *fakeCommand) ConnectIO() module.Command {
	c.ioConnected = true
	return c
}

func (c *fakeCommand) Exec() ([]byte, error) {
	return c.fakeExec()
}
