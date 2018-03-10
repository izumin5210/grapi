package template

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Command9acbc40e153d42b0b3e06ca83e7c53e3670d8cdc = "package main\n\nimport (\n\t\"os\"\n)\n\nfunc main() int {\n\tos.Exit(run())\n}\n\nfunc run() int {\n\treturn 0\n}\n"

// Command returns go-assets FileSystem
var Command = assets.NewFileSystem(map[string][]string{"/": []string{}, "/cmd": []string{}, "/cmd/{{ .name }}": []string{"run.go.tmpl"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519538948, 1519538948000000000),
		Data:     nil,
	}, "/cmd": &assets.File{
		Path:     "/cmd",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519539115, 1519539115000000000),
		Data:     nil,
	}, "/cmd/{{ .name }}": &assets.File{
		Path:     "/cmd/{{ .name }}",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519829745, 1519829745000000000),
		Data:     nil,
	}, "/cmd/{{ .name }}/run.go.tmpl": &assets.File{
		Path:     "/cmd/{{ .name }}/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519829745, 1519829745000000000),
		Data:     []byte(_Command9acbc40e153d42b0b3e06ca83e7c53e3670d8cdc),
	}}, "")
