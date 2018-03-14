package template

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Init8d21956ba8abe388f964e47be0f7e5d170a2fce5 = ""
var _Init23b808cac963edf44a497827f2a6eff5ddac970f = "required = [\n  \"github.com/golang/protobuf/protoc-gen-go\",\n  \"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway\",\n  \"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger\",\n]\n\n[prune]\n  go-tests = true\n  unused-packages = true\n\n  [[prune.project]]\n    name = \"github.com/grpc-ecosystem/grpc-gateway\"\n    non-go = false\n    unused-packages = false\n\n[[constraint]]\n{{- if .headUsed }}\n  branch = \"master\"\n{{- end }}\n  name = \"github.com/izumin5210/grapi\"\n{{- if not .headUsed }}\n  version = \"{{ .version }}\"\n{{- end }}\n"
var _Initd135936e91856b6159ac2eedcf89aa9f07773f82 = "package main\n\nimport (\n\t\"os\"\n\n\t\"{{ .importPath }}/app\"\n)\n\nfunc main() {\n\tos.Exit(run())\n}\n\nfunc run() int {\n\terr := app.Run()\n\tif err != nil {\n\t\treturn 1\n\t}\n\treturn 0\n}\n"
var _Initc051c9ff1a8e446bc9636d3144c2775a7e235322 = "[protoc]\nprotos_dir = \"./api/protos\"\nout_dir = \"./api\"\nimport_dirs = [\n  \"./vendor/github.com/grpc-ecosystem/grpc-gateway\",\n  \"./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis\",\n]\n\n  [[protoc.plugins]]\n  path = \"./vendor/github.com/golang/protobuf/protoc-gen-go\"\n  name = \"go\"\n  args = { plugins = \"grpc\" }\n\n  [[protoc.plugins]]\n  path = \"./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway\"\n  name = \"grpc-gateway\"\n  args = { logtostderr = true }\n\n  [[protoc.plugins]]\n  path = \"./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger\"\n  name = \"swagger\"\n  args = { logtostderr = true }\n"
var _Init38e76c5db8962fa825cf2bd8b23a2dc985c4513e = "*.so\n/vendor\n/bin\n/tmp\n"
var _Init71ed560e812a4261bc8b56d9feaef4800830e0b7 = ""
var _Initbc4053f4dd26ceb67e4646e8c1d2cc75897c4dd0 = "package app\n\nimport (\n\t\"github.com/izumin5210/grapi/pkg/grapiserver\"\n)\n\n// Run starts the grapiserver.\nfunc Run() error {\n\treturn grapiserver.New().\n\t\tAddRegisterGrpcServerImplFuncs(\n\t\t// TODO\n\t\t).\n\t\tAddRegisterGatewayHandlerFuncs(\n\t\t// TODO\n\t\t).\n\t\tServe()\n}\n"

// Init returns go-assets FileSystem
var Init = assets.NewFileSystem(map[string][]string{"/api": []string{}, "/api/protos": []string{".keep.tmpl"}, "/app": []string{"run.go.tmpl"}, "/app/server": []string{".keep.tmpl"}, "/cmd": []string{}, "/cmd/server": []string{"run.go.tmpl"}, "/": []string{".gitignore.tmpl", "Gopkg.toml.tmpl", "grapi.toml.tmpl"}}, map[string]*assets.File{
	"/app": &assets.File{
		Path:     "/app",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521038556, 1521038556000000000),
		Data:     nil,
	}, "/app/run.go.tmpl": &assets.File{
		Path:     "/app/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     []byte(_Initbc4053f4dd26ceb67e4646e8c1d2cc75897c4dd0),
	}, "/api/protos": &assets.File{
		Path:     "/api/protos",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/cmd/server": &assets.File{
		Path:     "/cmd/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/api/protos/.keep.tmpl": &assets.File{
		Path:     "/api/protos/.keep.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     []byte(_Init8d21956ba8abe388f964e47be0f7e5d170a2fce5),
	}, "/cmd": &assets.File{
		Path:     "/cmd",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/Gopkg.toml.tmpl": &assets.File{
		Path:     "/Gopkg.toml.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     []byte(_Init23b808cac963edf44a497827f2a6eff5ddac970f),
	}, "/.gitignore.tmpl": &assets.File{
		Path:     "/.gitignore.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     []byte(_Init38e76c5db8962fa825cf2bd8b23a2dc985c4513e),
	}, "/api": &assets.File{
		Path:     "/api",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/app/server": &assets.File{
		Path:     "/app/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521038585, 1521038585000000000),
		Data:     nil,
	}, "/app/server/.keep.tmpl": &assets.File{
		Path:     "/app/server/.keep.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1521038585, 1521038585000000000),
		Data:     []byte(_Init71ed560e812a4261bc8b56d9feaef4800830e0b7),
	}, "/cmd/server/run.go.tmpl": &assets.File{
		Path:     "/cmd/server/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     []byte(_Initd135936e91856b6159ac2eedcf89aa9f07773f82),
	}, "/grapi.toml.tmpl": &assets.File{
		Path:     "/grapi.toml.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     []byte(_Initc051c9ff1a8e446bc9636d3144c2775a7e235322),
	}}, "")
