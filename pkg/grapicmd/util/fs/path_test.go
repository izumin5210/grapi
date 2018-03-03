package fs

import (
	"testing"

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
