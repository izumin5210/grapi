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

	err := generator.GenerateService("foo/bar-baz")

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	files := []string{
		"api/protos/foo/bar_baz.proto",
		"app/server/foo/bar_baz_server.go",
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
}
