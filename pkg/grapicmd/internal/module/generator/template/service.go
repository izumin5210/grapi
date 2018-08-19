package template

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Servicedd0ea9770a89c1038381555c2b8b1dfc7d52b6d8 = "package {{.Go.Package }}\n\nimport (\n\t\"context\"\n{{range .Go.Imports}}\n\t\"{{.}}\"\n{{- end}}\n\n\t{{.PbGo.PackageName}} \"{{ .PbGo.PackagePath }}\"\n)\n\n// New{{.Go.ServerName}} creates a new {{.Go.ServerName}} instance.\nfunc New{{.Go.ServerName}}() interface {\n\t{{.PbGo.PackageName }}.{{.Go.ServerName}}\n\tgrapiserver.Server\n} {\n\treturn &{{.Go.StructName}}{}\n}\n\ntype {{.Go.StructName}} struct {\n}\n{{$go := .Go -}}\n{{$pbGo := .PbGo -}}\n{{- range .Methods}}\nfunc (s *{{$go.StructName}}) {{.Method}}(ctx context.Context, req *{{.RequestGo $pbGo.PackageName}}) (*{{.ResponseGo $pbGo.PackageName}}, error) {\n\t// TODO: Not yet implemented.\n\treturn nil, status.Error(codes.Unimplemented, \"TODO: You should implement it!\")\n}\n{{end -}}\n"
var _Servicea66d22b3414614e986072628e857757760b5fab0 = "package {{.Go.Package }}\n\nimport (\n\t\"context\"\n\n\t\"github.com/grpc-ecosystem/grpc-gateway/runtime\"\n\t\"google.golang.org/grpc\"\n\n\t{{.PbGo.PackageName}} \"{{ .PbGo.PackagePath }}\"\n)\n\n// RegisterWithServer implements grapiserver.Server.RegisterWithServer.\nfunc (s *{{.Go.StructName}}) RegisterWithServer(grpcSvr *grpc.Server) {\n\t{{.PbGo.PackageName}}.Register{{.Go.ServerName}}(grpcSvr, s)\n}\n\n// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.\nfunc (s *{{.Go.StructName}}) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {\n\treturn {{.PbGo.PackageName}}.Register{{.ServiceName}}ServiceHandler(ctx, mux, conn)\n}\n"
var _Service92aa9b7adaba175be79051b307a723ea2c536222 = "package {{.Go.Package }}\n{{if .Methods}}\nimport (\n\t\"context\"\n\t\"testing\"\n{{range .Go.TestImports}}\n\t\"{{.}}\"\n{{- end}}\n\n\t{{.PbGo.PackageName}} \"{{ .PbGo.PackagePath }}\"\n)\n{{$go := .Go -}}\n{{$pbGo := .PbGo -}}\n{{- range .Methods}}\nfunc Test_{{$go.ServerName}}_{{.Method}}(t *testing.T) {\n\tsvr := New{{$go.ServerName}}()\n\n\tctx := context.Background()\n\treq := &{{.RequestGo $pbGo.PackageName}}{}\n\n\tresp, err := svr.{{.Method}}(ctx, req)\n\n\tif err != nil {\n\t\tt.Errorf(\"returned an error %v\", err)\n\t}\n\n\tif resp == nil {\n\t\tt.Error(\"response should not nil\")\n\t}\n}\n{{end -}}\n{{end -}}\n"
var _Servicec080f048193e3b40b184b8e68c773e1b1bf56088 = "syntax = \"proto3\";\noption go_package = \"{{ .PbGo.PackagePath }};{{ .PbGo.PackageName }}\";\npackage {{ .Proto.Package }};\n{{range .Proto.Imports}}\nimport \"{{.}}\";\n{{- end}}\n\nservice {{ .ServiceName }}Service {\t\n{{- range .Methods}}\n  rpc {{.Method}} ({{.RequestProto}}) returns ({{.ResponseProto}}) {\n    option (google.api.http) = {\n      {{.HTTP.Method}}: \"/{{.HTTP.Path}}\"\n      {{- if .HTTP.Body}}\n      body: \"{{.HTTP.Body}}\"\n      {{- end}}\n    };\n  }\n{{- end}}\n}\n{{range .Proto.Messages}}\nmessage {{.Name}} {\n  {{- range .Fields}}\n  {{- if .Repeated}}\n  repeated {{.Type}} {{.Name}} = {{.Tag}};\n  {{- else}}\n  {{.Type}} {{.Name}} = {{.Tag}};\n  {{- end}}\n  {{- end}}\n}\n{{end -}}\n"

// Service returns go-assets FileSystem
var Service = assets.NewFileSystem(map[string][]string{"/": []string{}, "/{{.ProtoDir}}": []string{"{{.Path}}.proto.tmpl"}, "/{{.ServerDir}}": []string{"{{.Path}}_server.go.tmpl", "{{.Path}}_server_register_funcs.go.tmpl", "{{.Path}}_server_test.go.tmpl"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1530436336, 1530436336000000000),
		Data:     nil,
	}, "/{{.ProtoDir}}": &assets.File{
		Path:     "/{{.ProtoDir}}",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1534681932, 1534681932000000000),
		Data:     nil,
	}, "/{{.ProtoDir}}/{{.Path}}.proto.tmpl": &assets.File{
		Path:     "/{{.ProtoDir}}/{{.Path}}.proto.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1534681932, 1534681932000000000),
		Data:     []byte(_Servicec080f048193e3b40b184b8e68c773e1b1bf56088),
	}, "/{{.ServerDir}}": &assets.File{
		Path:     "/{{.ServerDir}}",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1531386686, 1531386686000000000),
		Data:     nil,
	}, "/{{.ServerDir}}/{{.Path}}_server.go.tmpl": &assets.File{
		Path:     "/{{.ServerDir}}/{{.Path}}_server.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1530436336, 1530436336000000000),
		Data:     []byte(_Servicedd0ea9770a89c1038381555c2b8b1dfc7d52b6d8),
	}, "/{{.ServerDir}}/{{.Path}}_server_register_funcs.go.tmpl": &assets.File{
		Path:     "/{{.ServerDir}}/{{.Path}}_server_register_funcs.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1530436336, 1530436336000000000),
		Data:     []byte(_Servicea66d22b3414614e986072628e857757760b5fab0),
	}, "/{{.ServerDir}}/{{.Path}}_server_test.go.tmpl": &assets.File{
		Path:     "/{{.ServerDir}}/{{.Path}}_server_test.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1531386686, 1531386686000000000),
		Data:     []byte(_Service92aa9b7adaba175be79051b307a723ea2c536222),
	}}, "")
