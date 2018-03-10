package template

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Command9acbc40e153d42b0b3e06ca83e7c53e3670d8cdc = "package main\n\nimport (\n\t\"fmt\"\n\t\"os\"\n)\n\nfunc main() {\n\tos.Exit(run())\n}\n\nfunc run() int {\n\tfmt.Println(\"It works!\")\n\treturn 0\n}\n"

// Command returns go-assets FileSystem
var Command = assets.NewFileSystem(map[string][]string{"/cmd": []string{}, "/cmd/{{ .name }}": []string{"run.go.tmpl"}, "/": []string{}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520680416, 1520680416000000000),
		Data:     nil,
	}, "/cmd": &assets.File{
		Path:     "/cmd",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520680416, 1520680416000000000),
		Data:     nil,
	}, "/cmd/{{ .name }}": &assets.File{
		Path:     "/cmd/{{ .name }}",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520680724, 1520680724000000000),
		Data:     nil,
	}, "/cmd/{{ .name }}/run.go.tmpl": &assets.File{
		Path:     "/cmd/{{ .name }}/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520680724, 1520680724000000000),
		Data:     []byte(_Command9acbc40e153d42b0b3e06ca83e7c53e3670d8cdc),
	}}, "")
