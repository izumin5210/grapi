package testing

import (
	"go/build"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type Ctx struct {
	GOPATH    string
	RootDir   cli.RootDir
	CreateCmd func(*testing.T, afero.Fs, Case) *cobra.Command
	Cases     []Case
}

type Case struct {
	Test         string
	Args         []string
	Files        []string
	SkippedFiles map[string]struct{}
	ProtoDir     string
	ProtoOutDir  string
	ServerDir    string
	PkgName      string
}

func Run(t *testing.T, ctx *Ctx) {
	t.Helper()

	defer func(c build.Context) { fs.BuildContext = c }(fs.BuildContext)
	fs.BuildContext = build.Context{GOPATH: ctx.GOPATH}

	for _, tc := range ctx.Cases {
		t.Run(tc.Test, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			cmd := ctx.CreateCmd(t, fs, tc)

			t.Run("generate", func(t *testing.T) {
				cmd.SetArgs(append([]string{"generate"}, tc.Args...))
				err := cmd.Execute()

				if err != nil {
					t.Errorf("returned an error: %+v", err)
				}

				for _, file := range tc.Files {
					t.Run(file, func(t *testing.T) {
						if _, ok := tc.SkippedFiles[file]; ok {
							ok, err := afero.Exists(fs, file)

							if err != nil {
								t.Errorf("returned an error: %v", err)
							}

							if ok {
								t.Error("should not exist")
							}
						} else {
							data, err := afero.ReadFile(fs, ctx.RootDir.Join(file))

							if err != nil {
								t.Errorf("returned an error: %v", err)
							}

							cupaloy.SnapshotT(t, string(data))
						}
					})
				}
			})

			t.Run("destroy", func(t *testing.T) {
				t.SkipNow()
				// err := generator.DestroyService(c.name)

				// if err != nil {
				// 	t.Errorf("returned an error: %v", err)
				// }

				// for _, file := range c.files {
				// 	t.Run(file, func(t *testing.T) {
				// 		ok, err := afero.Exists(fs, filepath.Join(rootDir, file))

				// 		if err != nil {
				// 			t.Errorf("Exists(fs, %q) returned an error: %v", file, err)
				// 		}

				// 		if ok {
				// 			t.Errorf("%q should not exist", file)
				// 		}
				// 	})
				// }
			})
		})
	}
}
