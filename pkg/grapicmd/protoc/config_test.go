package protoc

import (
	"path/filepath"
	"reflect"
	"testing"
)

func createDefaultConfig() *Config {
	return &Config{
		ImportDirs: []string{
			"./vendor/github.com/grpc-ecosystem/grpc-gateway",
			"./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis",
		},
		ProtosDir: "./api/protos",
		OutDir:    "./api",
		Plugins: []*Plugin{
			{
				Path: "./vendor/github.com/golang/protobuf/protoc-gen-go",
				Name: "go",
				Args: map[string]interface{}{"plugins": "grpc"},
			},
			{
				Path: "./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway",
				Name: "grpc-gateway",
				Args: map[string]interface{}{"logtostderr": true},
			},
			{
				Path: "./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger",
				Name: "swagger",
				Args: map[string]interface{}{"logtostderr": true},
			},
		},
	}
}

func Test_Config_OutDirOf(t *testing.T) {
	cfg := createDefaultConfig()
	rootDir := "/home/example/app"

	t.Run("with no errors", func(t *testing.T) {
		outDir, err := cfg.OutDirOf(rootDir, filepath.Join(rootDir, "api/protos/foo/bar.proto"))

		if err != nil {
			t.Errorf("OutDirOf returns an error %v", err)
		}

		if got, want := outDir, "api/foo"; got != want {
			t.Errorf("OutDirOf returned %q, want %q", got, want)
		}
	})

	t.Run("with an error", func(t *testing.T) {
		outDir, err := cfg.OutDirOf(rootDir, filepath.Join(rootDir, "api/foo/bar.proto"))

		if err == nil {
			t.Errorf("OutDirOf should return an error")
		}

		if got, want := outDir, ""; got != want {
			t.Errorf("OutDirOf returned %q, want %q", got, want)
		}
	})
}

func Test_Config_Commands(t *testing.T) {
	cfg := createDefaultConfig()
	rootDir := "/home/example/app"
	protoPath := "api/protos/foo/bar.proto"

	cmds, err := cfg.Commands(rootDir, filepath.Join(rootDir, protoPath))

	if err != nil {
		t.Errorf("Commands() returned an error %v", err)
	}

	wantCmd := []string{
		"protoc",
		"-I", "api/protos/foo",
		"-I", "./vendor/github.com/grpc-ecosystem/grpc-gateway",
		"-I", "./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis",
		"--go_out=plugins=grpc:api/foo",
		"api/protos/foo/bar.proto",
	}

	if got, want := len(cmds), len(cfg.Plugins); got != want {
		t.Errorf("Commands() returned %d commands, want %d commands", got, want)
	} else if got, want := cmds[0], wantCmd; !reflect.DeepEqual(got, want) {
		t.Errorf("Commands()[0] returned %v, want %v", got, want)
	}
}
