package grapicmd

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets47a88c8c8af6573dc96014ca4a3fe07cf60b2cf8 = "[prune]\n  go-tests = true\n  unused-packages = true\n\n[[constraint]]\n  branch = \"master\"\n  name = \"github.com/izumin5210/grapi\"\n"
var _Assets574eaaa50e5ef6944af024c75d238ab9966d7c8b = "package app\n\nimport (\n\t\"github.com/izumin5210/grapi/pkg/grapiserver\"\n)\n\n// Run starts the grapiserver.\nfunc Run() error {\n\treturn grapiserver.New().\n\t\tAddRegisterGrpcServerImplFuncs(\n\t\t// TODO\n\t\t).\n\t\tAddRegisterGatewayHandlerFuncs(\n\t\t// TODO\n\t\t).\n\t\tServe()\n}\n"
var _Assets0be0fbdbabaa4b7745f1449bcc16c946349450ca = ""
var _Assets4e67b53a3c4e87a4cdbcda5cb40ea75fa9b60dca = "package main\n\nimport (\n\t\"{{ .importPath }}/app\"\n)\n\n// Run starts the grapiserver.\nfunc Run(args []string) error {\n\treturn app.Run()\n}\n"
var _Assets3dfc78964ec5558e3912c9d65663751a137cf520 = "*.so\n/vendor\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/cmd/server": []string{"run.go.tmpl"}, "/": []string{".gitignore.tmpl", "Gopkg.toml.tmpl"}, "/api": []string{}, "/api/protos": []string{".keep.tmpl"}, "/app": []string{"run.go.tmpl"}, "/cmd": []string{}}, map[string]*assets.File{
	"/api": &assets.File{
		Path:     "/api",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519464132, 1519464132000000000),
		Data:     nil,
	}, "/app/run.go.tmpl": &assets.File{
		Path:     "/app/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519464304, 1519464304000000000),
		Data:     []byte(_Assets574eaaa50e5ef6944af024c75d238ab9966d7c8b),
	}, "/cmd/server": &assets.File{
		Path:     "/cmd/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519465416, 1519465416000000000),
		Data:     nil,
	}, "/Gopkg.toml.tmpl": &assets.File{
		Path:     "/Gopkg.toml.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519464377, 1519464377000000000),
		Data:     []byte(_Assets47a88c8c8af6573dc96014ca4a3fe07cf60b2cf8),
	}, "/cmd/server/run.go.tmpl": &assets.File{
		Path:     "/cmd/server/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519464297, 1519464297000000000),
		Data:     []byte(_Assets4e67b53a3c4e87a4cdbcda5cb40ea75fa9b60dca),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519467786, 1519467786000000000),
		Data:     nil,
	}, "/.gitignore.tmpl": &assets.File{
		Path:     "/.gitignore.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519464438, 1519464438000000000),
		Data:     []byte(_Assets3dfc78964ec5558e3912c9d65663751a137cf520),
	}, "/api/protos": &assets.File{
		Path:     "/api/protos",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519465398, 1519465398000000000),
		Data:     nil,
	}, "/api/protos/.keep.tmpl": &assets.File{
		Path:     "/api/protos/.keep.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519464155, 1519464155000000000),
		Data:     []byte(_Assets0be0fbdbabaa4b7745f1449bcc16c946349450ca),
	}, "/app": &assets.File{
		Path:     "/app",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519465406, 1519465406000000000),
		Data:     nil,
	}, "/cmd": &assets.File{
		Path:     "/cmd",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519464132, 1519464132000000000),
		Data:     nil,
	}}, "")
