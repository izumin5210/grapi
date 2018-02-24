package cmd

import (
	"bytes"
	"go/build"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
)

var (
	tmplPaths []string
)

func init() {
	tmplPaths = make([]string, 0, len(grapicmd.Assets.Files))
	for path := range grapicmd.Assets.Files {
		tmplPaths = append(tmplPaths, path)
	}
	sort.Strings(tmplPaths)
}

func newInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a grapi application",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to get the current working directory")
			}

			root := wd

			if argCnt := len(args); argCnt == 1 {
				arg := args[0]
				if arg != "." {
					if filepath.IsAbs(arg) {
						root = arg
					} else {
						root, err = filepath.Abs(arg)
						if err != nil {
							return errors.Wrap(err, "failed to get the target directory")
						}
					}
				}
			} else if argCnt > 1 {
				return errors.Errorf("invalid argument count: want 0 or 1, got %d", argCnt)
			}

			return errors.WithStack(initProject(afero.NewOsFs(), root))
		},
	}
}

func initProject(fs afero.Fs, rootPath string) error {
	var importPath string
	for _, gopath := range filepath.SplitList(build.Default.GOPATH) {
		prefix := filepath.Join(gopath, "src") + "/"
		// FIXME: should not use strings.HasPrefix
		if strings.HasPrefix(rootPath, prefix) {
			importPath = strings.Replace(rootPath, prefix, "", 1)
			break
		}
	}
	if importPath == "" {
		return errors.New("failed to get the import path")
	}
	data := map[string]string{
		"importPath": importPath,
	}
	for _, tmplPath := range tmplPaths {
		entry := grapicmd.Assets.Files[tmplPath]
		if entry.IsDir() {
			continue
		}
		path := strings.TrimSuffix(tmplPath, ".tmpl")
		absPath := filepath.Join(rootPath, path)
		dirPath := filepath.Dir(absPath)
		if ok, err := afero.DirExists(fs, dirPath); err != nil {
			return errors.Wrapf(err, "failed to retrieve %s", dirPath)
		} else if !ok {
			err = fs.MkdirAll(dirPath, 0755)
			if err != nil {
				return errors.Wrapf(err, "failed to create %s", dirPath)
			}
		}
		tmpl, err := template.New("").Parse(string(entry.Data))
		if err != nil {
			return errors.Wrapf(err, "failed to parse the template of %s", path)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, data)
		if err != nil {
			return errors.Wrapf(err, "failed to generate %s", path)
		}
		err = afero.WriteFile(fs, absPath, buf.Bytes(), 0644)
		if err != nil {
			return errors.Wrapf(err, "failed to write %s", path)
		}
	}
	return nil
}
