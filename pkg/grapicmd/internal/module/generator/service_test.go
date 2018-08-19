package generator

import (
	"go/build"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
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
	ui.EXPECT().ItemSkipped(gomock.Any()).AnyTimes()
	fs := afero.NewMemMapFs()

	cases := []struct {
		name         string
		args         []string
		files        []string
		skippedFiles map[string]struct{}
		scaffold     bool
		skipTest     bool
		resource     string
		protoDir     string
		protoOutDir  string
		serverDir    string
		pkgName      string
	}{
		{
			name: "foo",
			files: []string{
				"api/protos/foo.proto",
				"app/server/foo_server.go",
				"app/server/foo_server_register_funcs.go",
				"app/server/foo_server_test.go",
			},
		},
		{
			name: "foo",
			files: []string{
				"api/protos/foo.proto",
				"app/server/foo_server.go",
				"app/server/foo_server_register_funcs.go",
				"app/server/foo_server_test.go",
			},
			pkgName: "testcompany.testapp",
		},
		{
			name: "foo/bar",
			files: []string{
				"api/protos/foo/bar.proto",
				"app/server/foo/bar_server.go",
				"app/server/foo/bar_server_register_funcs.go",
				"app/server/foo/bar_server_test.go",
			},
		},
		{
			name: "foo/bar_baz",
			files: []string{
				"api/protos/foo/bar_baz.proto",
				"app/server/foo/bar_baz_server.go",
				"app/server/foo/bar_baz_server_register_funcs.go",
				"app/server/foo/bar_baz_server_test.go",
			},
		},
		{
			name: "foo/bar-baz",
			files: []string{
				"api/protos/foo/bar_baz.proto",
				"app/server/foo/bar_baz_server.go",
				"app/server/foo/bar_baz_server_register_funcs.go",
				"app/server/foo/bar_baz_server_test.go",
			},
		},
		{
			name: "foo/bar-baz",
			args: []string{"list", "create", "delete"},
			files: []string{
				"api/protos/foo/bar_baz.proto",
				"app/server/foo/bar_baz_server.go",
				"app/server/foo/bar_baz_server_register_funcs.go",
				"app/server/foo/bar_baz_server_test.go",
			},
		},
		{
			name: "foo/bar-baz",
			args: []string{"list", "create", "rename", "delete", "move_move"},
			files: []string{
				"api/protos/foo/bar_baz.proto",
				"app/server/foo/bar_baz_server.go",
				"app/server/foo/bar_baz_server_register_funcs.go",
				"app/server/foo/bar_baz_server_test.go",
			},
		},
		{
			name: "qux",
			files: []string{
				"pkg/foo/protos/qux.proto",
				"app/server/qux_server.go",
				"app/server/qux_server_register_funcs.go",
				"app/server/qux_server_test.go",
			},
			protoDir: "pkg/foo/protos",
		},
		{
			name: "quux",
			files: []string{
				"api/protos/quux.proto",
				"app/server/quux_server.go",
				"app/server/quux_server_register_funcs.go",
				"app/server/quux_server_test.go",
			},
			protoOutDir: "api/out",
		},
		{
			name: "corge",
			files: []string{
				"api/protos/corge.proto",
				"pkg/foo/server/corge_server.go",
				"pkg/foo/server/corge_server_register_funcs.go",
				"pkg/foo/server/corge_server_test.go",
			},
			serverDir: "pkg/foo/server",
		},
		{
			name: "book",
			files: []string{
				"api/protos/book.proto",
				"app/server/book_server.go",
				"app/server/book_server_register_funcs.go",
				"app/server/book_server_test.go",
			},
			scaffold: true,
		},
		{
			name: "book",
			files: []string{
				"api/protos/book.proto",
				"app/server/book_server.go",
				"app/server/book_server_register_funcs.go",
			},
			skippedFiles: map[string]struct{}{
				"app/server/book_server_test.go": {},
			},
			skipTest: true,
		},
		{
			name: "library",
			files: []string{
				"api/protos/library.proto",
				"app/server/library_server.go",
				"app/server/library_server_register_funcs.go",
				"app/server/library_server_test.go",
			},
			resource: "book",
			scaffold: true,
		},
	}

	for _, c := range cases {
		test := c.name
		if len(c.args) > 0 {
			test += " with " + strings.Join(c.args, ",")
		}

		generator := newServiceGenerator(fs, ui, rootDir, c.protoDir, c.protoOutDir, c.serverDir, c.pkgName)

		t.Run(test, func(t *testing.T) {
			test := "Generate"
			if c.scaffold {
				test = "Scaffold"
			}
			if c.skipTest {
				test += " without test"
			}
			t.Run(test, func(t *testing.T) {
				var err error
				if c.scaffold {
					err = generator.ScaffoldService(c.name, module.ServiceGenerationConfig{ResourceName: c.resource, SkipTest: c.skipTest})
				} else {
					err = generator.GenerateService(c.name, module.ServiceGenerationConfig{ResourceName: c.resource, Methods: c.args, SkipTest: c.skipTest})
				}

				if err != nil {
					t.Errorf("returned an error: %v", err)
				}

				for _, file := range c.files {
					t.Run(file, func(t *testing.T) {
						if _, ok := c.skippedFiles[file]; ok {
							ok, err := afero.Exists(fs, file)

							if err != nil {
								t.Errorf("returned an error: %v", err)
							}

							if ok {
								t.Error("should not exist")
							}
						} else {
							data, err := afero.ReadFile(fs, filepath.Join(rootDir, file))

							if err != nil {
								t.Errorf("returned an error: %v", err)
							}

							cupaloy.SnapshotT(t, string(data))
						}
					})
				}
			})

			t.Run("Destroy", func(t *testing.T) {
				err := generator.DestroyService(c.name)

				if err != nil {
					t.Errorf("returned an error: %v", err)
				}

				for _, file := range c.files {
					t.Run(file, func(t *testing.T) {
						ok, err := afero.Exists(fs, filepath.Join(rootDir, file))

						if err != nil {
							t.Errorf("Exists(fs, %q) returned an error: %v", file, err)
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
