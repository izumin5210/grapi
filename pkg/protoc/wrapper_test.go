package protoc_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/izumin5210/gex/pkg/tool"
	"github.com/spf13/afero"
	"k8s.io/utils/exec"
	"k8s.io/utils/exec/testing"

	"github.com/izumin5210/grapi/pkg/clui"
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

	execer := &testingexec.FakeExec{CommandScript: make([]testingexec.FakeCommandAction, 6)}
	for i := range execer.CommandScript {
		execer.CommandScript[i] = func(cmd string, args ...string) exec.Cmd {
			calls = append(calls, append([]string{cmd}, args...))
			return testingexec.InitFakeCmd(&testingexec.FakeCmd{
				CombinedOutputScript: []testingexec.FakeCombinedOutputAction{
					func() ([]byte, error) { return []byte("\n"), nil },
				},
			}, cmd, args...)
		}
	}

	rootDir := "/go/src/awesomeapp"
	binDir := filepath.Join(rootDir, "bin")
	protosDir := filepath.Join(rootDir, "api", "protos")

	fs := afero.NewMemMapFs()
	dieIf(t, fs.MkdirAll(binDir, 0755))
	dieIf(t, fs.MkdirAll(protosDir, 0755))
	dieIf(t, afero.WriteFile(fs, filepath.Join(rootDir, "api", "should_be_ignored.proto"), []byte{}, 0644))
	dieIf(t, afero.WriteFile(fs, filepath.Join(protosDir, "book.proto"), []byte{}, 0644))
	dieIf(t, afero.WriteFile(fs, filepath.Join(protosDir, "types", "users.proto"), []byte{}, 0644))

	cfg := &protoc.Config{
		ImportDirs: []string{
			filepath.Join(rootDir, "vendor", "github.com", "grpc-ecosystem", "grpc-gateway"),
			protosDir,
		},
		ProtosDir: "./api/protos",
		OutDir:    "./api",
		Plugins: []*protoc.Plugin{
			{Name: "go", Args: map[string]interface{}{"plugins": "grpc"}},
			{Name: "grpc-gateway", Args: map[string]interface{}{"logtostderr": true}},
			{Name: "swagger", Args: map[string]interface{}{"logtostderr": true}},
		},
	}

	wrapper := protoc.NewWrapper(cfg, fs, execer, clui.Nop, &fakeToolRepository{}, rootDir, binDir)

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
