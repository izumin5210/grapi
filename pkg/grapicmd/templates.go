package grapicmd

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets574eaaa50e5ef6944af024c75d238ab9966d7c8b = "package app\n\nimport (\n\t\"github.com/izumin5210/grapi/pkg/grapiserver\"\n)\n\n// Run starts the grapiserver.\nfunc Run() error {\n\treturn grapiserver.New().\n\t\tAddRegisterGrpcServerImplFuncs(\n\t\t\t// TODO\n\t\t).\n\t\tAddRegisterGatewayHandlerFuncs(\n\t\t\t// TODO\n\t\t).\n\t\tServe()\n}\n"
var _Assets4e67b53a3c4e87a4cdbcda5cb40ea75fa9b60dca = "package main\n\nimport (\n\t\"{{ .importPath }}/app\"\n)\n\n// Run starts the grapiserver.\nfunc Run(args []string) error {\n\treturn app.Run()\n}\n"
var _Assets3dfc78964ec5558e3912c9d65663751a137cf520 = "*.so\n/vendor\n/bin\n/tmp\n"
var _Assets0be0fbdbabaa4b7745f1449bcc16c946349450ca = ""
var _Assets47a88c8c8af6573dc96014ca4a3fe07cf60b2cf8 = "required = [\n  \"github.com/golang/protobuf/protoc-gen-go\",\n  \"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway\",\n  \"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger\",\n]\n\n[prune]\n  go-tests = true\n  unused-packages = true\n\n[[constraint]]\n  branch = \"master\"\n  name = \"github.com/izumin5210/grapi\"\n"
var _Assets41bfcc578717f32548542fb0270687ba4aa622dc = "[protoc]\nprotos_dir = \"./api/protos\"\nout_dir = \"./api\"\nimport_dirs = [\n  \"./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis\",\n]\n\n  [[protoc.plugins]]\n  name = \"go\"\n  args = { plugins = \"grpc\" }\n\n  [[protoc.plugins]]\n  name = \"grpc-gateway\"\n  args = { logtostderr = true }\n\n  [[protoc.plugins]]\n  name = \"swagger\"\n  args = { logtostderr = true }\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/app": []string{"run.go.tmpl"}, "/cmd": []string{}, "/cmd/server": []string{"run.go.tmpl"}, "/": []string{".gitignore.tmpl", "Gopkg.toml.tmpl", "grapi.toml.tmpl"}, "/api": []string{}, "/api/protos": []string{".keep.tmpl"}}, map[string]*assets.File{
	"/app/run.go.tmpl": &assets.File{
		Path:     "/app/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519480310, 1519480310000000000),
		Data:     []byte(_Assets574eaaa50e5ef6944af024c75d238ab9966d7c8b),
	}, "/cmd": &assets.File{
		Path:     "/cmd",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519464132, 1519464132000000000),
		Data:     nil,
	}, "/cmd/server": &assets.File{
		Path:     "/cmd/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519465416, 1519465416000000000),
		Data:     nil,
	}, "/cmd/server/run.go.tmpl": &assets.File{
		Path:     "/cmd/server/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519464297, 1519464297000000000),
		Data:     []byte(_Assets4e67b53a3c4e87a4cdbcda5cb40ea75fa9b60dca),
	}, "/api": &assets.File{
		Path:     "/api",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519464132, 1519464132000000000),
		Data:     nil,
	}, "/api/protos": &assets.File{
		Path:     "/api/protos",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519465398, 1519465398000000000),
		Data:     nil,
	}, "/.gitignore.tmpl": &assets.File{
		Path:     "/.gitignore.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519480280, 1519480280000000000),
		Data:     []byte(_Assets3dfc78964ec5558e3912c9d65663751a137cf520),
	}, "/api/protos/.keep.tmpl": &assets.File{
		Path:     "/api/protos/.keep.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519464155, 1519464155000000000),
		Data:     []byte(_Assets0be0fbdbabaa4b7745f1449bcc16c946349450ca),
	}, "/app": &assets.File{
		Path:     "/app",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519480310, 1519480310000000000),
		Data:     nil,
	}, "/Gopkg.toml.tmpl": &assets.File{
		Path:     "/Gopkg.toml.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519480975, 1519480975000000000),
		Data:     []byte(_Assets47a88c8c8af6573dc96014ca4a3fe07cf60b2cf8),
	}, "/grapi.toml.tmpl": &assets.File{
		Path:     "/grapi.toml.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519491180, 1519491180000000000),
		Data:     []byte(_Assets41bfcc578717f32548542fb0270687ba4aa622dc),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519491180, 1519491180000000000),
		Data:     nil,
	}}, "")
