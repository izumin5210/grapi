package protoc_test

import (
	"context"
	"os/exec"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/izumin5210/clig/pkg/clib"
	"github.com/izumin5210/execx"
	"github.com/izumin5210/gex/pkg/tool"
	"github.com/spf13/afero"

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

	exec := execx.New(execx.WithFakeProcess(
		func(ctx context.Context, cmd *exec.Cmd) error {
			switch cmd.Args[0] {
			case "go":
				if out := cmd.Stdout; out != nil {
					_, err := out.Write([]byte("/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.8.5\n"))
					dieIf(t, err)
				}
			case "protoc":
				calls = append(calls, cmd.Args)
			default:
				t.Errorf("unexpected command execution: %v", cmd.Args)
			}
			if out := cmd.Stdout; out != nil {
				_, err := out.Write([]byte("\n"))
				dieIf(t, err)
			}
			return nil
		},
	))

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

	wrapper := protoc.NewWrapper(cfg, fs, exec, cli.NopUI, &fakeToolRepository{}, rootDir)

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
