package fs

import (
	"errors"
	"go/build"
	"path/filepath"
	"strings"
)

func GetImportPath(rootPath string) (importPath string, err error) {
	for _, gopath := range filepath.SplitList(build.Default.GOPATH) {
		prefix := filepath.Join(gopath, "src") + "/"
		// FIXME: should not use strings.HasPrefix
		if strings.HasPrefix(rootPath, prefix) {
			importPath = strings.Replace(rootPath, prefix, "", 1)
			break
		}
	}
	if importPath == "" {
		err = errors.New("failed to get the import path")
	}
	return
}
