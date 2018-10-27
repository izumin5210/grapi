package fs

import (
	"os/user"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

func Test_GetImportPath(t *testing.T) {
	defer func(tmp string) { BuildContext.GOPATH = tmp }(BuildContext.GOPATH)
	BuildContext.GOPATH = "/home/go"

	cases := []struct {
		test  string
		in    string
		out   string
		isErr bool
	}{
		{
			test: "inside of GOPATH",
			in:   "/home/go/src/github.com/izumin5210/testapp",
			out:  "github.com/izumin5210/testapp",
		},
		{
			test: "directly under GOPATH",
			in:   "/home/go/src/testapp",
			out:  "testapp",
		},
		{
			test:  "outside of GOPATH",
			in:    "/home/go/testapp",
			isErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.test, func(t *testing.T) {
			out, err := GetImportPath(c.in)

			if got, want := err != nil, c.isErr; got != want {
				t.Errorf("Returned error is %t, want %t (%v)", got, want, err)
			}

			if got, want := out, c.out; got != want {
				t.Errorf("Returned %s, want %s", got, want)
			}
		})
	}
}

func Test_GetPackageName(t *testing.T) {
	defer func(tmp getOSUserFunc) { GetOSUser = tmp }(GetOSUser)
	GetOSUser = func() (*user.User, error) { return &user.User{Username: "testuser"}, nil }
	defer func(tmp string) { BuildContext.GOPATH = tmp }(BuildContext.GOPATH)
	BuildContext.GOPATH = "/home/go"

	cases := []struct {
		test  string
		in    string
		out   string
		isErr bool
	}{
		{
			test: "inside of GOPATH",
			in:   "/home/go/src/github.com/izumin5210/testapp",
			out:  "izumin5210.testapp",
		},
		{
			test: "directly under GOPATH",
			in:   "/home/go/src/testapp",
			out:  "testuser.testapp",
		},
		{
			test: "company name includes separators",
			in:   "/home/go/src/go.example.com/testapp",
			out:  "com.example.go.testapp",
		},
		{
			test: "package name includes hyphens",
			in:   "/home/go/src/go.example.com/test-app",
			out:  "com.example.go.test_app",
		},
		{
			test:  "outside of GOPATH",
			in:    "/home/go/testapp",
			isErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.test, func(t *testing.T) {
			out, err := GetPackageName(c.in)

			if got, want := err != nil, c.isErr; got != want {
				t.Errorf("Returned error is %t, want %t (%v)", got, want, err)
			}

			if got, want := out, c.out; got != want {
				t.Errorf("Returned %s, want %s", got, want)
			}
		})
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
