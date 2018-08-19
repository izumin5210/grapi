package template

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Initbc4053f4dd26ceb67e4646e8c1d2cc75897c4dd0 = "package app\n\nimport (\n\t\"github.com/izumin5210/grapi/pkg/grapiserver\"\n)\n\n// Run starts the grapiserver.\nfunc Run() error {\n\ts := grapiserver.New(\n\t\tgrapiserver.WithDefaultLogger(),\n\t\tgrapiserver.WithServers(\n\t\t// TODO\n\t\t),\n\t)\n\treturn s.Serve()\n}\n"
var _Init23b808cac963edf44a497827f2a6eff5ddac970f = "required = [\n  \"github.com/golang/protobuf/protoc-gen-go\",\n  \"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway\",\n  \"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger\",\n]\n\n[[constraint]]\n{{- if .headUsed }}\n  branch = \"master\"\n{{- end }}\n  name = \"github.com/izumin5210/grapi\"\n{{- if not .headUsed }}\n  version = \"{{ .version }}\"\n{{- end }}\n\n[[override]]\n  name = \"github.com/grpc-ecosystem/grpc-gateway\"\n  # To use `paths` option in `protoc-gen-grpc-gateway`\n  # See https://github.com/grpc-ecosystem/grpc-gateway/pull/711\n  revision = \"42fa2028fd3ba4a546dd525d9db71573f1992ba1\"\n"
var _Initc051c9ff1a8e446bc9636d3144c2775a7e235322 = "package = \"{{.packageName}}\"\n\n[grapi]\nserver_dir = \"./app/server\"\n\n[protoc]\nprotos_dir = \"./api/protos\"\nout_dir = \"./api\"\nimport_dirs = [\n  \"./vendor/github.com/grpc-ecosystem/grpc-gateway\",\n  \"./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis\",\n]\n\n  [[protoc.plugins]]\n  path = \"./vendor/github.com/golang/protobuf/protoc-gen-go\"\n  name = \"go\"\n  args = { plugins = \"grpc\", paths = \"source_relative\" }\n\n  [[protoc.plugins]]\n  path = \"./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway\"\n  name = \"grpc-gateway\"\n  args = { logtostderr = true, paths = \"source_relative\" }\n\n  [[protoc.plugins]]\n  path = \"./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger\"\n  name = \"swagger\"\n  args = { logtostderr = true }\n"
var _Initd135936e91856b6159ac2eedcf89aa9f07773f82 = "package main\n\nimport (\n\t\"os\"\n\n\t\"google.golang.org/grpc/grpclog\"\n\n\t\"{{ .importPath }}/app\"\n)\n\nfunc main() {\n\tos.Exit(run())\n}\n\nfunc run() int {\n\terr := app.Run()\n\tif err != nil {\n\t\tgrpclog.Errorf(\"server was shutdown with errors: %v\", err)\n\t\treturn 1\n\t}\n\treturn 0\n}\n"
var _Init38e76c5db8962fa825cf2bd8b23a2dc985c4513e = "*.so\n/vendor\n/bin\n/tmp\n"
var _Init8d21956ba8abe388f964e47be0f7e5d170a2fce5 = ""
var _Init71ed560e812a4261bc8b56d9feaef4800830e0b7 = ""

// Init returns go-assets FileSystem
var Init = assets.NewFileSystem(map[string][]string{"/cmd/server": []string{"run.go.tmpl"}, "/": []string{".gitignore.tmpl", "Gopkg.toml.tmpl", "grapi.toml.tmpl"}, "/api": []string{}, "/api/protos": []string{".keep.tmpl"}, "/app": []string{"run.go.tmpl"}, "/app/server": []string{".keep.tmpl"}, "/cmd": []string{}}, map[string]*assets.File{
	"/api/protos/.keep.tmpl": &assets.File{
		Path:     "/api/protos/.keep.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     []byte(_Init8d21956ba8abe388f964e47be0f7e5d170a2fce5),
	}, "/app/server/.keep.tmpl": &assets.File{
		Path:     "/app/server/.keep.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1521042119, 1521042119000000000),
		Data:     []byte(_Init71ed560e812a4261bc8b56d9feaef4800830e0b7),
	}, "/cmd/server": &assets.File{
		Path:     "/cmd/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1530424640, 1530424640000000000),
		Data:     nil,
	}, "/grapi.toml.tmpl": &assets.File{
		Path:     "/grapi.toml.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1534650220, 1534650220000000000),
		Data:     []byte(_Initc051c9ff1a8e446bc9636d3144c2775a7e235322),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1534650220, 1534650220000000000),
		Data:     nil,
	}, "/app/run.go.tmpl": &assets.File{
		Path:     "/app/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1521995193, 1521995193000000000),
		Data:     []byte(_Initbc4053f4dd26ceb67e4646e8c1d2cc75897c4dd0),
	}, "/app/server": &assets.File{
		Path:     "/app/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521042119, 1521042119000000000),
		Data:     nil,
	}, "/Gopkg.toml.tmpl": &assets.File{
		Path:     "/Gopkg.toml.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1534649282, 1534649282000000000),
		Data:     []byte(_Init23b808cac963edf44a497827f2a6eff5ddac970f),
	}, "/api": &assets.File{
		Path:     "/api",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/api/protos": &assets.File{
		Path:     "/api/protos",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/cmd": &assets.File{
		Path:     "/cmd",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/cmd/server/run.go.tmpl": &assets.File{
		Path:     "/cmd/server/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1530424640, 1530424640000000000),
		Data:     []byte(_Initd135936e91856b6159ac2eedcf89aa9f07773f82),
	}, "/.gitignore.tmpl": &assets.File{
		Path:     "/.gitignore.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     []byte(_Init38e76c5db8962fa825cf2bd8b23a2dc985c4513e),
	}, "/app": &assets.File{
		Path:     "/app",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521995193, 1521995193000000000),
		Data:     nil,
	}}, "")
