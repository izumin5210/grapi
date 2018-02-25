package cmd

import (
	"bytes"
	"fmt"
	"go/build"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/izumin5210/grapi/pkg/grapicmd"
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

func newInitCommand(out io.Writer) *cobra.Command {
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

			return errors.WithStack(initProject(afero.NewOsFs(), out, root))
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
	colorAttrsByStatus = map[status][]color.Attribute{
		statusCreate:     {color.FgGreen, color.Bold},
		statusExist:      {color.FgBlue, color.Bold},
		statusIdentical:  {color.FgBlue, color.Bold},
		statusConflicted: {color.FgRed, color.Bold},
		statusForce:      {color.FgYellow, color.Bold},
		statusSkipped:    {color.FgYellow, color.Bold},
	}
	creatableStatusSet = map[status]struct{}{
		statusCreate: struct{}{},
		statusForce:  struct{}{},
	}
)

func (s status) String() string {
	return nameByStatus[s]
}

func (s status) Fprintln(out io.Writer, msg string) {
	colored := color.New(colorAttrsByStatus[s]...).SprintfFunc()
	fmt.Fprintf(out, "%s  %s\n", colored("%12s", s.String()), msg)
}

func (s status) ShouldCreate() bool {
	_, ok := creatableStatusSet[s]
	return ok
}

func initProject(afs afero.Fs, out io.Writer, rootPath string) error {
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

		st.Fprintln(out, path[1:])
	}
	return nil
}
