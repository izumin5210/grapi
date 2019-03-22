package protoc

import (
	"path/filepath"
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
				Name: "go",
				Args: map[string]interface{}{"plugins": "grpc"},
			},
			{
				Name: "grpc-gateway",
				Args: map[string]interface{}{"logtostderr": true},
			},
			{
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
