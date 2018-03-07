package internal

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/izumin5210/grapi/pkg/grapicmd/command"
	"github.com/spf13/afero"
)

func Test_Script(t *testing.T) {
	type testContext struct {
		factory    ScriptFactory
		fs         afero.Fs
		rootDir    string
		executions [][]string
	}

	createTestContext := func(t *testing.T) *testContext {
		ctx := &testContext{
			fs:         afero.NewMemMapFs(),
			rootDir:    "/home/app",
			executions: [][]string{},
		}

		executor := &fakeExecutor{
			fakeExec: func(nameAndArgs []string, _ ...command.Option) ([]byte, error) {
				ctx.executions = append(ctx.executions, nameAndArgs)
				return []byte{}, nil
			},
		}

		ctx.factory = NewScriptFactory(ctx.fs, executor, ctx.rootDir)

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

	if got, want := testCtx.executions[0], []string{"go", "build", "-v", "-o=" + bin, entry}; !reflect.DeepEqual(got, want) {
		t.Errorf("Build() executed %v, want %v", got, want)
	}

	err = s.Run()
	if err != nil {
		t.Errorf("Run() returned an error %v", err)
	}

	if got, want := testCtx.executions[1], []string{bin}; !reflect.DeepEqual(got, want) {
		t.Errorf("Run() executed %v, want %v", got, want)
	}
}

// fake impls

type fakeExecutor struct {
	command.Executor
	fakeExec func([]string, ...command.Option) ([]byte, error)
}

func (e *fakeExecutor) Exec(nameAndArgs []string, opts ...command.Option) ([]byte, error) {
	return e.fakeExec(nameAndArgs, opts...)
}
