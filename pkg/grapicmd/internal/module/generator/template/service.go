package template

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Serviceba3c18adf9d77d5cbf8a40461aa014bc49f81edc = "syntax = \"proto3\";\noption go_package = \"{{ .PbGo.PackageName }}\";\npackage {{ .Proto.Package }};\n{{range .Proto.Imports}}\nimport \"{{.}}\";\n{{- end}}\n\nservice {{ .ServiceName }}Service {\t\n{{- range .Methods}}\n  rpc {{.Method}} ({{.RequestProto}}) returns ({{.ResponseProto}}) {\n    option (google.api.http) = {\n      {{.HTTP.Method}}: \"/{{.HTTP.Path}}\"\n      {{- if .HTTP.Body}}\n      body: \"{{.HTTP.Body}}\"\n      {{- end}}\n    };\n  }\n{{- end}}\n}\n{{range .Proto.Messages}}\nmessage {{.Name}} {\n  {{- range .Fields}}\n  {{- if .Repeated}}\n  repeated {{.Type}} {{.Name}} = {{.Tag}};\n  {{- else}}\n  {{.Type}} {{.Name}} = {{.Tag}};\n  {{- end}}\n  {{- end}}\n}\n{{end -}}\n"
var _Serviceef91c225ded973f86ca6b58050abcc766e9d41c3 = "package {{.Go.Package }}\n\nimport (\n\t\"context\"\n{{range .Go.Imports}}\n\t\"{{.}}\"\n{{- end}}\n\n\t{{.PbGo.PackageName}} \"{{ .PbGo.PackagePath }}\"\n)\n\n// New{{.Go.ServerName}} creates a new {{.Go.ServerName}} instance.\nfunc New{{.Go.ServerName}}() interface {\n\t{{.PbGo.PackageName }}.{{.Go.ServerName}}\n\tgrapiserver.Server\n} {\n\treturn &{{.Go.StructName}}{}\n}\n\ntype {{.Go.StructName}} struct {\n}\n\n// RegisterWithServer implements grapiserver.Server.RegisterWithServer.\nfunc (s *{{.Go.StructName}}) RegisterWithServer(grpcSvr *grpc.Server) {\n\t{{.PbGo.PackageName}}.Register{{.Go.ServerName}}(grpcSvr, s)\n}\n\n// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.\nfunc (s *{{.Go.StructName}}) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {\n\treturn {{.PbGo.PackageName}}.Register{{.ServiceName}}ServiceHandler(ctx, mux, conn)\n}\n{{$go := .Go -}}\n{{$pbGo := .PbGo -}}\n{{- range .Methods}}\nfunc (s *{{$go.StructName}}) {{.Method}}(ctx context.Context, req *{{.RequestGo $pbGo.PackageName}}) (*{{.ResponseGo $pbGo.PackageName}}, error) {\n\t// TODO: Not yet implemented.\n\treturn nil, status.Error(codes.Unimplemented, \"TODO: You should implement it!\")\n}\n{{end -}}\n"
var _Service7664b60536a30580d1374729ea1d6df5fd874e1e = "package {{.Go.Package }}\n{{if .Methods}}\nimport (\n\t\"context\"\n\t\"testing\"\n{{range .Go.TestImports}}\n\t\"{{.}}\"\n{{- end}}\n\n\t{{.PbGo.PackageName}} \"{{ .PbGo.PackagePath }}\"\n)\n{{$go := .Go -}}\n{{$pbGo := .PbGo -}}\n{{- range .Methods}}\nfunc Test_{{$go.ServerName}}_{{.Method}}(t *testing.T) {\n\tsvr := New{{$go.ServerName}}()\n\n\tctx := context.Background()\n\treq := &{{.RequestGo $pbGo.PackageName}}{}\n\n\tresp, err := svr.{{.Method}}(ctx, req)\n\n\tif err != nil {\n\t\tt.Errorf(\"returned an error %v\", err)\n\t}\n\n\tif resp == nil {\n\t\tat.Error(\"response should not nil\")\n\t}\n}\n{{end -}}\n{{end -}}\n"

// Service returns go-assets FileSystem
var Service = assets.NewFileSystem(map[string][]string{"/app": []string{}, "/app/server": []string{"{{.Path}}_server.go.tmpl", "{{.Path}}_server_test.go.tmpl"}, "/": []string{}, "/api": []string{}, "/api/protos": []string{"{{.Path}}.proto.tmpl"}}, map[string]*assets.File{
	"/api": &assets.File{
		Path:     "/api",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521913317, 1521913317000000000),
		Data:     nil,
	}, "/api/protos": &assets.File{
		Path:     "/api/protos",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521913317, 1521913317000000000),
		Data:     nil,
	}, "/api/protos/{{.Path}}.proto.tmpl": &assets.File{
		Path:     "/api/protos/{{.Path}}.proto.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1521913317, 1521913317000000000),
		Data:     []byte(_Serviceba3c18adf9d77d5cbf8a40461aa014bc49f81edc),
	}, "/app/server/{{.Path}}_server.go.tmpl": &assets.File{
		Path:     "/app/server/{{.Path}}_server.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1521995496, 1521995496000000000),
		Data:     []byte(_Serviceef91c225ded973f86ca6b58050abcc766e9d41c3),
	}, "/app/server/{{.Path}}_server_test.go.tmpl": &assets.File{
		Path:     "/app/server/{{.Path}}_server_test.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1521995220, 1521995220000000000),
		Data:     []byte(_Service7664b60536a30580d1374729ea1d6df5fd874e1e),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521913317, 1521913317000000000),
		Data:     nil,
	}, "/app": &assets.File{
		Path:     "/app",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521913317, 1521913317000000000),
		Data:     nil,
	}, "/app/server": &assets.File{
		Path:     "/app/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521995496, 1521995496000000000),
		Data:     nil,
	}}, "")
