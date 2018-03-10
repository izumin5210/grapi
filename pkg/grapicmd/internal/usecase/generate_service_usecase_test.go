package usecase

import (
	"go/build"
	"testing"

	"github.com/golang/mock/gomock"

	moduletesting "github.com/izumin5210/grapi/pkg/grapicmd/internal/module/testing"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

func Test_GenerateServiceUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tmpBuildContext := fs.BuildContext
	defer func() { fs.BuildContext = tmpBuildContext }()
	fs.BuildContext = build.Context{
		GOPATH: "/home",
	}

	rootDir := "/home/src/foo"

	type Case struct {
		input           string
		importPath      string
		path            string
		name            string
		serviceName     string
		packagePath     string
		packageName     string
		pbgoPackagePath string
		pbgoPackageName string
		protoPackage    string
	}

	cases := []Case{
		{
			input:           "bar",
			importPath:      "foo",
			path:            "bar",
			name:            "bar",
			serviceName:     "Bar",
			packagePath:     "server",
			packageName:     "server",
			pbgoPackagePath: "api",
			pbgoPackageName: "api_pb",
			protoPackage:    "foo.api",
		},
		{
			input:           "bar/baz",
			importPath:      "foo",
			path:            "bar/baz",
			name:            "baz",
			serviceName:     "Baz",
			packagePath:     "bar",
			packageName:     "bar",
			pbgoPackagePath: "api/bar",
			pbgoPackageName: "bar_pb",
			protoPackage:    "foo.api.bar",
		},
		{
			input:           "bar/baz/qux",
			importPath:      "foo",
			path:            "bar/baz/qux",
			name:            "qux",
			serviceName:     "Qux",
			packagePath:     "bar/baz",
			packageName:     "baz",
			pbgoPackagePath: "api/bar/baz",
			pbgoPackageName: "baz_pb",
			protoPackage:    "foo.api.bar.baz",
		},
		{
			input:           "bar/baz/qux_quux",
			importPath:      "foo",
			path:            "bar/baz/qux_quux",
			name:            "qux_quux",
			serviceName:     "QuxQuux",
			packagePath:     "bar/baz",
			packageName:     "baz",
			pbgoPackagePath: "api/bar/baz",
			pbgoPackageName: "baz_pb",
			protoPackage:    "foo.api.bar.baz",
		},
		{
			input:           "bar/baz/qux-quux",
			importPath:      "foo",
			path:            "bar/baz/qux_quux",
			name:            "qux_quux",
			serviceName:     "QuxQuux",
			packagePath:     "bar/baz",
			packageName:     "baz",
			pbgoPackagePath: "api/bar/baz",
			pbgoPackageName: "baz_pb",
			protoPackage:    "foo.api.bar.baz",
		},
		{
			input:           "bar-baz/qux-quux",
			importPath:      "foo",
			path:            "bar_baz/qux_quux",
			name:            "qux_quux",
			serviceName:     "QuxQuux",
			packagePath:     "bar_baz",
			packageName:     "bar_baz",
			pbgoPackagePath: "api/bar_baz",
			pbgoPackageName: "bar_baz_pb",
			protoPackage:    "foo.api.bar_baz",
		},
	}

	for _, c := range cases {
		ui := moduletesting.NewMockUI(ctrl)
		ui.EXPECT().Section(gomock.Any()).AnyTimes()
		ui.EXPECT().Subsection(gomock.Any()).AnyTimes()

		generator := moduletesting.NewMockGenerator(ctrl)
		generator.EXPECT().Exec(rootDir, map[string]interface{}{
			"importPath":      c.importPath,
			"path":            c.path,
			"name":            c.name,
			"serviceName":     c.serviceName,
			"packagePath":     c.packagePath,
			"packageName":     c.packageName,
			"pbgoPackagePath": c.pbgoPackagePath,
			"pbgoPackageName": c.pbgoPackageName,
			"protoPackage":    c.protoPackage,
		})

		usecase := NewGenerateServiceUsecase(ui, generator, rootDir)
		err := usecase.Perform(c.input)

		if err != nil {
			t.Errorf("Perform() returned an error %v", err)
		}
	}
}
