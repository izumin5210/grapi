package template

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Init38e76c5db8962fa825cf2bd8b23a2dc985c4513e = "*.so\n/vendor\n/bin\n/tmp\n"
var _Init50bb4ac2099b3758964058926b3c90524e478a2c = ""
var _Initce49182536bf5d7fb269f752e2e9d15dd5eaf7d6 = "package main\n\nimport (\n\t\"os\"\n\n\t\"google.golang.org/grpc/grpclog\"\n)\n\nfunc main() {\n\terr := run()\n\tif err != nil {\n\t\tgrpclog.Errorf(\"server was shutdown with errors: %v\", err)\n\t\tos.Exit(1)\n\t}\n}\n"
var _Init8d21956ba8abe388f964e47be0f7e5d170a2fce5 = ""
var _Initc051c9ff1a8e446bc9636d3144c2775a7e235322 = "package = \"{{.packageName}}\"\n\n[grapi]\nserver_dir = \"./app/server\"\n\n[protoc]\nprotos_dir = \"./api/protos\"\nout_dir = \"./api\"\nimport_dirs = [\n  \"./api/protos\",\n  '{{`{{ module \"github.com/grpc-ecosystem/grpc-gateway\" }}`}}',\n  '{{`{{ module \"github.com/grpc-ecosystem/grpc-gateway\" }}`}}/third_party/googleapis',\n]\n\n  [[protoc.plugins]]\n  name = \"go\"\n  args = { plugins = \"grpc\", paths = \"source_relative\" }\n\n  [[protoc.plugins]]\n  name = \"grpc-gateway\"\n  args = { logtostderr = true, paths = \"source_relative\" }\n\n  [[protoc.plugins]]\n  name = \"swagger\"\n  args = { logtostderr = true }\n"
var _Init71ed560e812a4261bc8b56d9feaef4800830e0b7 = ""
var _Initd135936e91856b6159ac2eedcf89aa9f07773f82 = "package main\n\nimport (\n\t\"github.com/izumin5210/grapi/pkg/grapiserver\"\n)\n\nfunc run() error {\n\ts := grapiserver.New(\n\t\tgrapiserver.WithDefaultLogger(),\n\t\tgrapiserver.WithServers(\n\t\t// TODO\n\t\t),\n\t)\n\treturn s.Serve()\n}\n"
var _Init881d845139e03b8e1e2dafd175e793d03f9bacaf = "// +build tools\n\npackage tools\n\n// tool dependencies\nimport (\n\t_ \"github.com/golang/protobuf/protoc-gen-go\"\n\t_ \"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway\"\n\t_ \"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger\"\n\t_ \"github.com/izumin5210/grapi/cmd/grapi\"\n\t_ \"github.com/izumin5210/grapi/cmd/grapi-gen-command\"\n\t_ \"github.com/izumin5210/grapi/cmd/grapi-gen-scaffold-service\"\n\t_ \"github.com/izumin5210/grapi/cmd/grapi-gen-service\"\n\t_ \"github.com/izumin5210/grapi/cmd/grapi-gen-type\"\n)\n"

// Init returns go-assets FileSystem
var Init = assets.NewFileSystem(map[string][]string{"/app/server": []string{".keep.tmpl"}, "/api/protos/type": []string{".keep.tmpl"}, "/api/protos": []string{".keep.tmpl"}, "/": []string{".gitignore.tmpl", "tools.go", "grapi.toml.tmpl"}, "/cmd": []string{}, "/cmd/server": []string{"main.go.tmpl", "run.go.tmpl"}, "/app": []string{}, "/api": []string{}}, map[string]*assets.File{
	"/cmd/server/run.go.tmpl": &assets.File{
		Path:     "/cmd/server/run.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1553312167, 1553312167203472863),
		Data:     []byte(_Initd135936e91856b6159ac2eedcf89aa9f07773f82),
	}, "/tools.go": &assets.File{
		Path:     "/tools.go",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1553227614, 1553227614146442611),
		Data:     []byte(_Init881d845139e03b8e1e2dafd175e793d03f9bacaf),
	}, "/api/protos/type/.keep.tmpl": &assets.File{
		Path:     "/api/protos/type/.keep.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1546318740, 1546318740832736215),
		Data:     []byte(_Init50bb4ac2099b3758964058926b3c90524e478a2c),
	}, "/cmd": &assets.File{
		Path:     "/cmd",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1546318740, 1546318740833267962),
		Data:     nil,
	}, "/.gitignore.tmpl": &assets.File{
		Path:     "/.gitignore.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1546318740, 1546318740832281063),
		Data:     []byte(_Init38e76c5db8962fa825cf2bd8b23a2dc985c4513e),
	}, "/api/protos": &assets.File{
		Path:     "/api/protos",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1546318740, 1546318740832681569),
		Data:     nil,
	}, "/api/protos/type": &assets.File{
		Path:     "/api/protos/type",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1546318740, 1546318740832742956),
		Data:     nil,
	}, "/api/protos/.keep.tmpl": &assets.File{
		Path:     "/api/protos/.keep.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1546318740, 1546318740832594918),
		Data:     []byte(_Init8d21956ba8abe388f964e47be0f7e5d170a2fce5),
	}, "/grapi.toml.tmpl": &assets.File{
		Path:     "/grapi.toml.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1553227614, 1553227614145264680),
		Data:     []byte(_Initc051c9ff1a8e446bc9636d3144c2775a7e235322),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1553227614, 1553227614145475871),
		Data:     nil,
	}, "/cmd/server/main.go.tmpl": &assets.File{
		Path:     "/cmd/server/main.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1553313981, 1553313981458672258),
		Data:     []byte(_Initce49182536bf5d7fb269f752e2e9d15dd5eaf7d6),
	}, "/app/server": &assets.File{
		Path:     "/app/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1546318740, 1546318740833088256),
		Data:     nil,
	}, "/api": &assets.File{
		Path:     "/api",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1546318740, 1546318740832528992),
		Data:     nil,
	}, "/cmd/server": &assets.File{
		Path:     "/cmd/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1553313981, 1553313981459028463),
		Data:     nil,
	}, "/app": &assets.File{
		Path:     "/app",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1553312311, 1553312311023436070),
		Data:     nil,
	}, "/app/server/.keep.tmpl": &assets.File{
		Path:     "/app/server/.keep.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1546318740, 1546318740833081345),
		Data:     []byte(_Init71ed560e812a4261bc8b56d9feaef4800830e0b7),
	}}, "")
