package main

import (
	"context"
	"go/build"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/gencmd"
	gencmdtesting "github.com/izumin5210/grapi/pkg/gencmd/testing"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
	"github.com/izumin5210/grapi/pkg/protoc"
	"github.com/izumin5210/grapi/pkg/svcgen"
	svcgentesting "github.com/izumin5210/grapi/pkg/svcgen/testing"
)

func TestRun(t *testing.T) {
	defer func(c build.Context) { fs.BuildContext = c }(fs.BuildContext)
	fs.BuildContext = build.Context{GOPATH: "/home"}

	rootDir := cli.RootDir("/home/src/testapp")

	cases := []struct {
		test         string
		args         []string
		files        []string
		skippedFiles map[string]struct{}
		protoDir     string
		protoOutDir  string
		serverDir    string
		pkgName      string
	}{
		{
			test: "simple",
			args: []string{"foo"},
			files: []string{
				"api/protos/foo.proto",
				"app/server/foo_server.go",
				"app/server/foo_server_register_funcs.go",
				"app/server/foo_server_test.go",
			},
		},
		{
			test: "specify package",
			args: []string{"foo"},
			files: []string{
				"api/protos/foo.proto",
				"app/server/foo_server.go",
				"app/server/foo_server_register_funcs.go",
				"app/server/foo_server_test.go",
			},
			pkgName: "testcompany.testapp",
		},
		{
			test: "nested",
			args: []string{"foo/bar"},
			files: []string{
				"api/protos/foo/bar.proto",
				"app/server/foo/bar_server.go",
				"app/server/foo/bar_server_register_funcs.go",
				"app/server/foo/bar_server_test.go",
			},
		},
		{
			test: "nested with specify pacakge",
			args: []string{"foo/bar"},
			files: []string{
				"api/protos/foo/bar.proto",
				"app/server/foo/bar_server.go",
				"app/server/foo/bar_server_register_funcs.go",
				"app/server/foo/bar_server_test.go",
			},
			pkgName: "testcompany.testapp",
		},
		{
			test: "snake_case name",
			args: []string{"foo/bar_baz"},
			files: []string{
				"api/protos/foo/bar_baz.proto",
				"app/server/foo/bar_baz_server.go",
				"app/server/foo/bar_baz_server_register_funcs.go",
				"app/server/foo/bar_baz_server_test.go",
			},
		},
		{
			test: "kebab-case name",
			args: []string{"foo/bar-baz"},
			files: []string{
				"api/protos/foo/bar_baz.proto",
				"app/server/foo/bar_baz_server.go",
				"app/server/foo/bar_baz_server_register_funcs.go",
				"app/server/foo/bar_baz_server_test.go",
			},
		},
		{
			test: "with some standard methods",
			args: []string{"foo/bar-baz", "list", "create", "delete"},
			files: []string{
				"api/protos/foo/bar_baz.proto",
				"app/server/foo/bar_baz_server.go",
				"app/server/foo/bar_baz_server_register_funcs.go",
				"app/server/foo/bar_baz_server_test.go",
			},
		},
		{
			test: "with non-standard methods",
			args: []string{"foo/bar-baz", "list", "create", "rename", "delete", "move_move"},
			files: []string{
				"api/protos/foo/bar_baz.proto",
				"app/server/foo/bar_baz_server.go",
				"app/server/foo/bar_baz_server_register_funcs.go",
				"app/server/foo/bar_baz_server_test.go",
			},
		},
		{
			test: "specify proto dir",
			args: []string{"qux"},
			files: []string{
				"pkg/foo/protos/qux.proto",
				"app/server/qux_server.go",
				"app/server/qux_server_register_funcs.go",
				"app/server/qux_server_test.go",
			},
			protoDir: "pkg/foo/protos",
		},
		{
			test: "specify proto out dir",
			args: []string{"quux"},
			files: []string{
				"api/protos/quux.proto",
				"app/server/quux_server.go",
				"app/server/quux_server_register_funcs.go",
				"app/server/quux_server_test.go",
			},
			protoOutDir: "api/out",
		},
		{
			test: "specify server dir",
			args: []string{"corge"},
			files: []string{
				"api/protos/corge.proto",
				"pkg/foo/server/corge_server.go",
				"pkg/foo/server/corge_server_register_funcs.go",
				"pkg/foo/server/corge_server_test.go",
			},
			serverDir: "pkg/foo/server",
		},
		{
			test: "skip tests",
			args: []string{"--skip-test", "book"},
			files: []string{
				"api/protos/book.proto",
				"app/server/book_server.go",
				"app/server/book_server_register_funcs.go",
			},
			skippedFiles: map[string]struct{}{
				"app/server/book_server_test.go": {},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.test, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			createSvcApp := func(ctx *gencmd.Ctx, cmd *gencmd.Command) (*svcgen.App, error) {
				return svcgentesting.NewTestApp(ctx, cmd, &fakeProtocWrapper{}, cli.NopUI)
			}
			createGenApp := func(ctx *gencmd.Ctx, cmd *gencmd.Command) (*gencmd.App, error) {
				return gencmdtesting.NewTestApp(ctx, cmd, cli.NopUI)
			}

			ctx := &grapicmd.Ctx{
				FS:      fs,
				RootDir: rootDir,
				Config: grapicmd.Config{
					Package: c.pkgName,
				},
				ProtocConfig: protoc.Config{
					ProtosDir: c.protoDir,
					OutDir:    c.protoOutDir,
				},
			}
			ctx.Config.Grapi.ServerDir = c.serverDir
			cmd := gencmd.NewCommand("service", &gencmd.Ctx{
				Ctx:           ctx,
				CreateAppFunc: createGenApp,
				GenerateCmd:   NewGenerateCommand(createSvcApp),
			})

			t.Run("generate", func(t *testing.T) {
				cmd.SetArgs(append([]string{"generate"}, c.args...))
				err := cmd.Execute()

				if err != nil {
					t.Errorf("returned an error: %+v", err)
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
							data, err := afero.ReadFile(fs, rootDir.Join(file))

							if err != nil {
								t.Errorf("returned an error: %v", err)
							}

							cupaloy.SnapshotT(t, string(data))
						}
					})
				}
			})

			// t.Run("Destroy", func(t *testing.T) {
			// 	err := generator.DestroyService(c.name)

			// 	if err != nil {
			// 		t.Errorf("returned an error: %v", err)
			// 	}

			// 	for _, file := range c.files {
			// 		t.Run(file, func(t *testing.T) {
			// 			ok, err := afero.Exists(fs, filepath.Join(rootDir, file))

			// 			if err != nil {
			// 				t.Errorf("Exists(fs, %q) returned an error: %v", file, err)
			// 			}

			// 			if ok {
			// 				t.Errorf("%q should not exist", file)
			// 			}
			// 		})
			// 	}
			// })
		})
	}
}

type fakeProtocWrapper struct{}

func (*fakeProtocWrapper) Exec(context.Context) error { return nil }
