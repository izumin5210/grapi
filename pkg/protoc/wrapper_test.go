package protoc_test

import (
	"context"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/izumin5210/clig/pkg/clib"
	"github.com/izumin5210/gex/pkg/tool"
	"github.com/spf13/afero"
	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/protoc"
)

func TestWrapper_Exec(t *testing.T) {
	dieIf := func(t *testing.T, err error) {
		t.Helper()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	calls := [][]string{}

	execer := &testingexec.FakeExec{
		CommandScript: make([]testingexec.FakeCommandAction, 7),
	}
	for i := range execer.CommandScript {
		execer.CommandScript[i] = func(cmd string, args ...string) exec.Cmd {
			switch cmd {
			case "go":
				out := []byte("/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.8.5\n")
				return testingexec.InitFakeCmd(&testingexec.FakeCmd{
					RunScript: []testingexec.FakeRunAction{
						func() ([]byte, []byte, error) { return out, nil, nil },
					},
				}, cmd, args...)
			case "protoc":
				calls = append(calls, append([]string{cmd}, args...))
			default:
				t.Errorf("unexpected command execution: %s %v", cmd, args)
			}
			return testingexec.InitFakeCmd(&testingexec.FakeCmd{
				CombinedOutputScript: []testingexec.FakeCombinedOutputAction{
					func() ([]byte, error) { return []byte("\n"), nil },
				},
			}, cmd, args...)
		}
	}

	rootDir := cli.RootDir{Path: clib.Path("/go/src/awesomeapp")}
	protosDir := rootDir.Join("api", "protos")

	fs := afero.NewMemMapFs()
	dieIf(t, fs.MkdirAll(rootDir.BinDir().String(), 0755))
	dieIf(t, fs.MkdirAll(protosDir.String(), 0755))
	dieIf(t, afero.WriteFile(fs, rootDir.Join("go.mod").String(), []byte("module example.com/awesomeapp"), 0644))
	dieIf(t, afero.WriteFile(fs, rootDir.Join("go.sum").String(), []byte(""), 0644))
	dieIf(t, afero.WriteFile(fs, rootDir.Join("api", "should_be_ignored.proto").String(), []byte{}, 0644))
	dieIf(t, afero.WriteFile(fs, rootDir.Join("api", "should_be_ignored_proto").String(), []byte{}, 0644))
	dieIf(t, afero.WriteFile(fs, protosDir.Join("book.proto").String(), []byte{}, 0644))
	dieIf(t, afero.WriteFile(fs, protosDir.Join("types", "users.proto").String(), []byte{}, 0644))

	cfg := &protoc.Config{
		ImportDirs: []string{
			`{{ module "github.com/grpc-ecosystem/grpc-gateway" }}`,
			`{{ module "github.com/grpc-ecosystem/grpc-gateway" }}/third_party/googleapis`,
			protosDir.String(),
		},
		ProtosDir: "./api/protos",
		OutDir:    "./api",
		Plugins: []*protoc.Plugin{
			{Name: "go", Args: map[string]interface{}{"plugins": "grpc"}},
			{Name: "grpc-gateway", Args: map[string]interface{}{"logtostderr": true}},
			{Name: "swagger", Args: map[string]interface{}{"logtostderr": true}},
		},
	}

	wrapper := protoc.NewWrapper(cfg, fs, execer, cli.NopUI, &fakeToolRepository{}, rootDir)

	err := wrapper.Exec(context.TODO())
	if err != nil {
		t.Fatalf("returned %v, want nil", err)
	}

	cupaloy.SnapshotT(t, calls)
}

type fakeToolRepository struct {
	tool.Repository
}

func (*fakeToolRepository) BuildAll(context.Context) error { return nil }
