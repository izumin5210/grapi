package generator

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"

	moduletesting "github.com/izumin5210/grapi/pkg/grapicmd/internal/module/testing"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

func Test_ProjectGenerator(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func(tmp string) { fs.BuildContext.GOPATH = tmp }(fs.BuildContext.GOPATH)
	fs.BuildContext.GOPATH = "/home"

	rootDir := "/home/src/testcompany/testapp"

	ui := moduletesting.NewMockUI(ctrl)
	ui.EXPECT().ItemSuccess(gomock.Any()).AnyTimes()
	fs := afero.NewMemMapFs()

	generator := newProjectGenerator(fs, ui, "")

	err := generator.GenerateProject(rootDir, "", false)

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	files := []string{
		".gitignore",
		"Gopkg.toml",
		"grapi.toml",
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
}
