package fs

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

func Test_LookupRoot(t *testing.T) {
	fs := afero.NewMemMapFs()

	fs.MkdirAll("/home/root/app/server", 0755)
	fs.MkdirAll("/home/root/api/proto", 0755)
	fs.MkdirAll("/home/root/cmd/server", 0755)
	fs.MkdirAll("/home/app", 0755)
	afero.WriteFile(fs, "/home/root/grapi.toml", []byte(""), 0644)
	afero.WriteFile(fs, "/home/app/grapi.toml", []byte(""), 0644)

	cases := []struct {
		cwd  string
		root string
		ok   bool
	}{
		{cwd: "/", root: "", ok: false},
		{cwd: "/home", root: "", ok: false},
		{cwd: "/home/root", root: "/home/root", ok: true},
		{cwd: "/home/root/app", root: "/home/root", ok: true},
		{cwd: "/home/root/api/proto", root: "/home/root", ok: true},
		{cwd: "/home/app", root: "", ok: false},
		{cwd: "/home/api", root: "", ok: false},
	}

	for _, c := range cases {
		root, ok := LookupRoot(fs, c.cwd)

		if got, want := ok, c.ok; got != want {
			t.Errorf("LookupRoot(fs, %q) returns %t, want %t", c.cwd, got, want)
		}

		if got, want := root, c.root; got != want {
			t.Errorf("LookupRoot(fs, %q) returns %q, want %q", c.cwd, got, want)
		}
	}
}

func Test_FindUserDefinedCommandPaths(t *testing.T) {
	fs := afero.NewMemMapFs()
	rootPath := "/home/app"
	cmdPath := filepath.Join(rootPath, "cmd")

	fs.MkdirAll(filepath.Join(rootPath, "app", "server"), 0755)
	fs.MkdirAll(filepath.Join(rootPath, "api", "proto"), 0755)
	afero.WriteFile(fs, filepath.Join(rootPath, "grapi.toml"), []byte(""), 0644)
	fs.MkdirAll(filepath.Join(cmdPath, "server"), 0755)
	afero.WriteFile(fs, filepath.Join(cmdPath, "server", "run.go"), []byte("package main"), 0644)
	fs.MkdirAll(filepath.Join(cmdPath, "foo", "bar"), 0755)
	afero.WriteFile(fs, filepath.Join(cmdPath, "foo", "bar", "run.go"), []byte("package main"), 0644)
	fs.MkdirAll(filepath.Join(cmdPath, "baz"), 0755)
	afero.WriteFile(fs, filepath.Join(cmdPath, "baz", "qux"), []byte(""), 0644)
	fs.MkdirAll(filepath.Join(cmdPath, "quux", "corge"), 0755)
	afero.WriteFile(fs, filepath.Join(cmdPath, "quux", "main.go"), []byte("package main"), 0644)
	afero.WriteFile(fs, filepath.Join(cmdPath, "quux", "corge", "grault.go"), []byte("package corge"), 0644)
	fs.MkdirAll(filepath.Join(cmdPath, "garply"), 0755)
	afero.WriteFile(fs, filepath.Join(cmdPath, "garply", "main.go"), []byte("package garply"), 0644)
	fs.MkdirAll(filepath.Join(cmdPath, "waldo", "fred"), 0755)
	afero.WriteFile(fs, filepath.Join(cmdPath, "waldo", "main.go"), []byte("package waldo"), 0644)
	afero.WriteFile(fs, filepath.Join(cmdPath, "waldo", "fred", "main.go"), []byte("package main"), 0644)
	afero.WriteFile(fs, filepath.Join(cmdPath, "waldo", "fred", "main_test.go"), []byte("package main"), 0644)
	afero.WriteFile(fs, filepath.Join(cmdPath, "waldo", "fred", "plugh.go"), []byte("package main"), 0644)

	paths, err := FindMainPackagesAndSources(fs, cmdPath)

	if err != nil {
		t.Errorf("FindUserDefinedCommandPaths returned an error %v", err)
	}

	wantPaths := map[string][]string{
		filepath.Join(cmdPath, "foo", "bar"):    {"run.go"},
		filepath.Join(cmdPath, "quux"):          {"main.go"},
		filepath.Join(cmdPath, "waldo", "fred"): {"main.go", "plugh.go"},
		filepath.Join(cmdPath, "server"):        {"run.go"},
	}

	if got, want := len(paths), len(wantPaths); got != want {
		t.Errorf("FindUserDefinedCommandPaths returned %d paths, want %d paths", got, want)
	}

	if diff := cmp.Diff(paths, wantPaths); diff != "" {
		t.Errorf("Received path differs: (-got +want)\n%s", diff)
	}
}
