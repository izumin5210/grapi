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

	cases := []struct {
		test, in, out string
		isErr         bool
	}{
		{
			test: "simple",
			in:   "api/protos/foo/bar.proto",
			out:  "api/foo_pb",
		},
		{
			test:  "out of proto dir",
			in:    "api/foo/bar.proto",
			isErr: true,
		},
		{
			test: "directly under proto dir",
			in:   "api/protos/baz.proto",
			out:  "api",
		},
	}

	for _, c := range cases {
		t.Run(c.test, func(t *testing.T) {
			outDir, err := cfg.OutDirOf(rootDir, filepath.Join(rootDir, c.in))

			if got, want := err != nil, c.isErr; got != want {
				t.Errorf("OutDirOf returned an error: got %t, want %t", got, want)
			}

			if got, want := outDir, c.out; got != want {
				t.Errorf("OutDirOf returned %q, want %q", got, want)
			}
		})
	}
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
		"--go_out=plugins=grpc:api/foo_pb",
		"api/protos/foo/bar.proto",
	}

	if got, want := len(cmds), len(cfg.Plugins); got != want {
		t.Errorf("Commands() returned %d commands, want %d commands", got, want)
	} else if got, want := cmds[0], wantCmd; !reflect.DeepEqual(got, want) {
		t.Errorf("Commands()[0] returned %v, want %v", got, want)
	}
}
