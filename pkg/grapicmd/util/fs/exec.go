package fs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/afero"
)

func ListExecutableWithPrefix(fs afero.Fs, prefix string) []string {
	var wg sync.WaitGroup
	ch := make(chan string)

	for _, path := range filepath.SplitList(os.Getenv("PATH")) {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			files, err := afero.ReadDir(fs, path)
			if err != nil {
				return
			}

			for _, f := range files {
				if m := f.Mode(); !f.IsDir() && m&0111 != 0 && strings.HasPrefix(f.Name(), prefix) {
					ch <- f.Name()
				}
			}
		}(path)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	execs := make([]string, 0, 100)
	for exec := range ch {
		execs = append(execs, exec)
	}
	sort.Strings(execs)

	return execs
}
