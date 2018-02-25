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
	"github.com/izumin5210/grapi/pkg/grapicmd/ui"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

var (
	tmplPaths []string
)

func init() {
	rootFiles := make([]string, 0, len(grapicmd.Assets.Files))
	tmplPaths = make([]string, 0, len(grapicmd.Assets.Files))
	for path, entry := range grapicmd.Assets.Files {
		if entry.IsDir() {
			continue
		}
		if strings.Count(entry.Path[1:], "/") == 0 {
			rootFiles = append(rootFiles, path)
		} else {
			tmplPaths = append(tmplPaths, path)
		}
	}
	sort.Strings(rootFiles)
	sort.Strings(tmplPaths)
	tmplPaths = append(rootFiles, tmplPaths...)
}

func newInitCommand(ui ui.UI) *cobra.Command {
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

			return errors.WithStack(initProject(afero.NewOsFs(), ui, root))
		},
	}
}

type status int

const (
	statusCreate status = iota
	statusExist
	statusIdentical
	statusConflicted
	statusForce
	statusSkipped
)

var (
	nameByStatus = map[status]string{
		statusCreate:     "create",
		statusExist:      "exist",
		statusIdentical:  "identical",
		statusConflicted: "conflicted",
		statusForce:      "force",
		statusSkipped:    "skipped",
	}
	levelByStatus = map[status]ui.Level{
		statusCreate:     ui.LevelSuccess,
		statusExist:      ui.LevelInfo,
		statusIdentical:  ui.LevelInfo,
		statusConflicted: ui.LevelFail,
		statusForce:      ui.LevelWarn,
		statusSkipped:    ui.LevelWarn,
	}
	creatableStatusSet = map[status]struct{}{
		statusCreate: struct{}{},
		statusForce:  struct{}{},
	}
)

func (s status) String() string {
	return nameByStatus[s]
}

func (s status) Level() ui.Level {
	return levelByStatus[s]
}

func (s status) ShouldCreate() bool {
	_, ok := creatableStatusSet[s]
	return ok
}

func initProject(afs afero.Fs, ui ui.UI, rootPath string) error {
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
		path := strings.TrimSuffix(tmplPath, ".tmpl")
		absPath := filepath.Join(rootPath, path)
		dirPath := filepath.Dir(absPath)

		// create directory if not exists
		if err := fs.CreateDirIfNotExists(afs, dirPath); err != nil {
			return errors.WithStack(err)
		}

		// generate content
		tmpl, err := template.New("").Parse(string(entry.Data))
		if err != nil {
			return errors.Wrapf(err, "failed to parse the template of %s", path)
		}
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, data)
		if err != nil {
			return errors.Wrapf(err, "failed to generate %s", path)
		}

		// check existed entries
		st := statusCreate
		if ok, err := afero.Exists(afs, path); err != nil {
			// TODO: handle an error
			st = statusSkipped
		} else if ok {
			body, err := afero.ReadFile(afs, path)
			if err != nil {
				// TODO: handle an error
				st = statusSkipped
			}
			if string(body) == buf.String() {
				st = statusIdentical
			} else {
				// TODO: ask to overwrite
				st = statusConflicted
			}
		}

		// create
		if st.ShouldCreate() {
			err = afero.WriteFile(afs, absPath, buf.Bytes(), 0644)
			if err != nil {
				return errors.Wrapf(err, "failed to write %s", path)
			}
		}

		ui.PrintWithStatus(path[1:], st)
	}
	return nil
}
