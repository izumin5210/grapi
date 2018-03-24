package generator

import (
	"go/build"
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"

	moduletesting "github.com/izumin5210/grapi/pkg/grapicmd/internal/module/testing"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

func Test_ServiceGenerator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmpBuildContext := fs.BuildContext
	defer func() { fs.BuildContext = tmpBuildContext }()
	fs.BuildContext = build.Context{
		GOPATH: "/home",
	}

	rootDir := "/home/src/testapp"

	ui := moduletesting.NewMockUI(ctrl)
	ui.EXPECT().ItemSuccess(gomock.Any()).AnyTimes()
	fs := afero.NewMemMapFs()

	generator := newServiceGenerator(fs, ui, rootDir)

	cases := []struct {
		name  string
		files []string
	}{
		{
			name: "foo",
			files: []string{
				"api/protos/foo.proto",
				"app/server/foo_server.go",
			},
		},
		{
			name: "foo/bar",
			files: []string{
				"api/protos/foo/bar.proto",
				"app/server/foo/bar_server.go",
			},
		},
		{
			name: "foo/bar_baz",
			files: []string{
				"api/protos/foo/bar_baz.proto",
				"app/server/foo/bar_baz_server.go",
			},
		},
		{
			name: "foo/bar-baz",
			files: []string{
				"api/protos/foo/bar_baz.proto",
				"app/server/foo/bar_baz_server.go",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Run("Generate", func(t *testing.T) {
				err := generator.GenerateService(c.name)

				if err != nil {
					t.Errorf("returned an error %v", err)
				}

				for _, file := range c.files {
					t.Run(file, func(t *testing.T) {
						data, err := afero.ReadFile(fs, filepath.Join(rootDir, file))

						if err != nil {
							t.Errorf("returned an error %v", err)
						}

						cupaloy.SnapshotT(t, string(data))
					})
				}
			})

			t.Run("Destroy", func(t *testing.T) {
				err := generator.DestroyService(c.name)

				if err != nil {
					t.Errorf("returned an error %v", err)
				}

				for _, file := range c.files {
					t.Run(file, func(t *testing.T) {
						ok, err := afero.Exists(fs, filepath.Join(rootDir, file))

						if err != nil {
							t.Errorf("Exists(fs, %q) returned an error %v", file, err)
						}

						if ok {
							t.Errorf("%q should not exist", file)
						}
					})
				}
			})
		})
	}
}
