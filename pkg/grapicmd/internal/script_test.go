package internal

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/spf13/afero"
)

func Test_Script(t *testing.T) {
	type execution struct {
		nameAndArgs []string
		dir         string
		connected   bool
	}
	type testContext struct {
		factory    ScriptFactory
		fs         afero.Fs
		rootDir    string
		executions []*execution
	}

	createTestContext := func(t *testing.T) *testContext {
		ctx := &testContext{
			fs:         afero.NewMemMapFs(),
			rootDir:    "/home/app",
			executions: []*execution{},
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

		ctx.factory = NewScriptFactory(ctx.fs, commandFactory, ctx.rootDir)

		return ctx
	}

	testCtx := createTestContext(t)

	name := "bar"
	entry := filepath.Join(testCtx.rootDir, "cmd", name, "run.go")
	bin := filepath.Join(testCtx.rootDir, "bin", name)
	s := testCtx.factory.Create(entry)

	if got, want := s.Name(), name; got != want {
		t.Errorf("Name() returned %v, want %v", got, want)
	}

	err := s.Build()
	if err != nil {
		t.Errorf("Build() returned an error %v", err)
	}

	exec := testCtx.executions[0]

	if got, want := exec.nameAndArgs, []string{"go", "build", "-v", "-o=" + bin, entry}; !reflect.DeepEqual(got, want) {
		t.Errorf("Build() executed %v, want %v", got, want)
	}

	if got, want := exec.dir, "/home/app"; got != want {
		t.Errorf("Build() executed a command in %v, want %v", got, want)
	}

	err = s.Run()
	if err != nil {
		t.Errorf("Run() returned an error %v", err)
	}

	exec = testCtx.executions[1]

	if got, want := exec.nameAndArgs, []string{bin}; !reflect.DeepEqual(got, want) {
		t.Errorf("Run() executed %v, want %v", got, want)
	}

	if got, want := exec.dir, "/home/app"; got != want {
		t.Errorf("Run() executed a command in %v, want %v", got, want)
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
