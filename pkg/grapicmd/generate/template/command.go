package template

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Command9acbc40e153d42b0b3e06ca83e7c53e3670d8cdc = "package main\n\n// Run starts awesome process.\nfunc Run(args []string) error {\n\treturn nil\n}\n"

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
		Mtime:    time.Unix(1519539347, 1519539347000000000),
		Data:     nil,
	}, "/cmd/{{ .name }}/run.go.tmpl": &assets.File{
		Path:     "/cmd/{{ .name }}/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519539045, 1519539045000000000),
		Data:     []byte(_Command9acbc40e153d42b0b3e06ca83e7c53e3670d8cdc),
	}}, "")
