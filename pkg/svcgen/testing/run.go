package testing

import (
	"go/build"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/gencmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
	"github.com/spf13/afero"
)

type Ctx struct {
	GOPATH    string
	RootDir   cli.RootDir
	CreateCmd func(*testing.T, afero.Fs, Case) gencmd.Executor
	Cases     []Case
}

type Case struct {
	Test         string
	GArgs        []string
	DArgs        []string
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
			afero.WriteFile(fs, ctx.RootDir.Join("grapi.toml").String(), []byte{}, 0755)

			t.Run("generate", func(t *testing.T) {
				cmd := ctx.CreateCmd(t, fs, tc)
				cmd.Command().SetArgs(append([]string{"generate"}, tc.GArgs...))
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
							data, err := afero.ReadFile(fs, ctx.RootDir.Join(file).String())

							if err != nil {
								t.Errorf("returned an error: %v", err)
							}

							cupaloy.SnapshotT(t, string(data))
						}
					})
				}
			})

			t.Run("destroy", func(t *testing.T) {
				cmd := ctx.CreateCmd(t, fs, tc)
				cmd.Command().SetArgs(append([]string{"destroy"}, tc.DArgs...))
				err := cmd.Execute()

				if err != nil {
					t.Errorf("returned an error: %+v", err)
				}

				for _, file := range tc.Files {
					t.Run(file, func(t *testing.T) {
						ok, err := afero.Exists(fs, ctx.RootDir.Join(file).String())

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
