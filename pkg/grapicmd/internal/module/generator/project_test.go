package generator

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

func TestProjectGenerator_GenerateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func(tmp string) { fs.BuildContext.GOPATH = tmp }(fs.BuildContext.GOPATH)
	fs.BuildContext.GOPATH = "/home"

	rootDir := "/home/src/testcompany/testapp"

	cases := []struct {
		test string
	}{
		{
			test: "no config",
		},
	}

	for _, tc := range cases {
		t.Run(tc.test, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			generator := newProjectGenerator(fs, cli.NopUI)
			err := generator.GenerateProject(rootDir, "")

			if err != nil {
				t.Errorf("returned an error %v", err)
			}

			files := []string{
				".gitignore",
				"grapi.toml",
				"tools.go",
				"app/run.go",
				"cmd/server/run.go",
			}

			for _, file := range files {
				t.Run(file, func(t *testing.T) {
					data, err := afero.ReadFile(fs, filepath.Join(rootDir, file))

					if err != nil {
						t.Errorf("returned an error %v", err)
					}

					cupaloy.SnapshotT(t, string(data))
				})
			}
		})
	}
}
