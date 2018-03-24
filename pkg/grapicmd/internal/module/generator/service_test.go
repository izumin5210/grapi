package generator

import (
	"go/build"
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"

	moduletesting "github.com/izumin5210/grapi/pkg/grapicmd/internal/module/testing"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

func Test_ServiceGenerator_Generator(t *testing.T) {
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

func Test_ServiceGenerator_createParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmpBuildContext := fs.BuildContext
	defer func() { fs.BuildContext = tmpBuildContext }()
	fs.BuildContext = build.Context{
		GOPATH: "/home",
	}

	rootDir := "/home/src/foo"

	type Case struct {
		input            string
		importPath       string
		path             string
		name             string
		serviceName      string
		localServiceName string
		packagePath      string
		packageName      string
		pbgoPackagePath  string
		pbgoPackageName  string
		protoPackage     string
	}

	cases := []Case{
		{
			input:            "bar",
			importPath:       "foo",
			path:             "bar",
			name:             "bar",
			serviceName:      "Bar",
			localServiceName: "bar",
			packagePath:      "server",
			packageName:      "server",
			pbgoPackagePath:  "api",
			pbgoPackageName:  "api_pb",
			protoPackage:     "foo.api",
		},
		{
			input:            "bar/baz",
			importPath:       "foo",
			path:             "bar/baz",
			name:             "baz",
			serviceName:      "Baz",
			localServiceName: "baz",
			packagePath:      "bar",
			packageName:      "bar",
			pbgoPackagePath:  "api/bar",
			pbgoPackageName:  "bar_pb",
			protoPackage:     "foo.api.bar",
		},
		{
			input:            "bar/baz/qux",
			importPath:       "foo",
			path:             "bar/baz/qux",
			name:             "qux",
			serviceName:      "Qux",
			localServiceName: "qux",
			packagePath:      "bar/baz",
			packageName:      "baz",
			pbgoPackagePath:  "api/bar/baz",
			pbgoPackageName:  "baz_pb",
			protoPackage:     "foo.api.bar.baz",
		},
		{
			input:            "bar/baz/qux_quux",
			importPath:       "foo",
			path:             "bar/baz/qux_quux",
			name:             "qux_quux",
			serviceName:      "QuxQuux",
			localServiceName: "quxQuux",
			packagePath:      "bar/baz",
			packageName:      "baz",
			pbgoPackagePath:  "api/bar/baz",
			pbgoPackageName:  "baz_pb",
			protoPackage:     "foo.api.bar.baz",
		},
		{
			input:            "bar/baz/qux-quux",
			importPath:       "foo",
			path:             "bar/baz/qux_quux",
			name:             "qux_quux",
			serviceName:      "QuxQuux",
			localServiceName: "quxQuux",
			packagePath:      "bar/baz",
			packageName:      "baz",
			pbgoPackagePath:  "api/bar/baz",
			pbgoPackageName:  "baz_pb",
			protoPackage:     "foo.api.bar.baz",
		},
		{
			input:            "bar-baz/qux-quux",
			importPath:       "foo",
			path:             "bar_baz/qux_quux",
			name:             "qux_quux",
			serviceName:      "QuxQuux",
			localServiceName: "quxQuux",
			packagePath:      "bar_baz",
			packageName:      "bar_baz",
			pbgoPackagePath:  "api/bar_baz",
			pbgoPackageName:  "bar_baz_pb",
			protoPackage:     "foo.api.bar_baz",
		},
	}

	for _, c := range cases {
		ui := moduletesting.NewMockUI(ctrl)
		fs := afero.NewMemMapFs()

		generator := newServiceGenerator(fs, ui, rootDir).(*serviceGenerator)

		got, err := generator.createParams(c.input)

		if err != nil {
			t.Errorf("Perform() returned an error %v", err)
		}

		want := map[string]interface{}{
			"importPath":       c.importPath,
			"path":             c.path,
			"name":             c.name,
			"serviceName":      c.serviceName,
			"localServiceName": c.localServiceName,
			"packagePath":      c.packagePath,
			"packageName":      c.packageName,
			"pbgoPackagePath":  c.pbgoPackagePath,
			"pbgoPackageName":  c.pbgoPackageName,
			"protoPackage":     c.protoPackage,
		}

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Received params differs: (-got +want)\n%s", diff)
		}
	}
}
