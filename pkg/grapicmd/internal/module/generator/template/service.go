package template

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Serviceba3c18adf9d77d5cbf8a40461aa014bc49f81edc = "syntax = \"proto3\";\noption go_package = \"{{ .PbGo.PackageName }}\";\npackage {{ .Proto.Package }};\n{{range .Proto.Imports}}\nimport \"{{.}}\";\n{{- end}}\n\nservice {{ .ServiceName }}Service {\t\n{{- range .Methods}}\n  rpc {{.Method}} ({{.RequestProto}}) returns ({{.ResponseProto}}) {\n    option (google.api.http) = {\n      {{.HTTP.Method}}: \"/{{.HTTP.Path}}\"\n      {{- if .HTTP.Body}}\n      body: \"{{.HTTP.Body}}\"\n      {{- end}}\n    };\n  }\n{{- end}}\n}\n{{range .Proto.Messages}}\nmessage {{.Name}} {\n  {{- range $i, $f := .Fields}}\n  {{- if .Repeated}}\n  repeated {{$f.Type}} {{$f.Name}} = {{$i}};\n  {{- else}}\n  {{$f.Type}} {{$f.Name}} = {{$i}};\n  {{- end}}\n  {{- end}}\n}\n{{end -}}\n"
var _Serviceef91c225ded973f86ca6b58050abcc766e9d41c3 = "package {{.Go.Package }}\n\nimport (\n\t\"context\"\n{{range .Go.Imports}}\n\t\"{{.}}\"\n{{- end}}\n\n\t{{.PbGo.PackageName}} \"{{ .PbGo.PackagePath }}\"\n)\n\nvar (\n\t// Register{{.ServiceName}}ServiceHandler is a function to register card service handler to gRPC Gateway's mux.\n\tRegister{{.ServiceName}}ServiceHandler = {{.PbGo.PackageName}}.Register{{.ServiceName}}ServiceHandler\n)\n\n// Register{{.Go.ServerName}}Factory creates a function to register card service server impl to grpc.Server.\nfunc Register{{.Go.ServerName}}Factory() func(s *grpc.Server) {\n\treturn func(s *grpc.Server) {\n\t\t{{.PbGo.PackageName}}.Register{{.Go.ServerName}}(s, New())\n\t}\n}\n\n// New creates a new {{.Go.ServerName}} instance.\nfunc New() {{.PbGo.PackageName }}.{{.Go.ServerName}} {\n\treturn &{{.Go.StructName}}{}\n}\n\ntype {{.Go.StructName}} struct {\n}\n{{$go := .Go -}}\n{{$pbGo := .PbGo -}}\n{{- range .Methods}}\nfunc (s *{{$go.StructName}}) {{.Method}}(ctx context.Context, req *{{.RequestGo $pbGo.PackageName}}) (*{{.ResponseGo $pbGo.PackageName}}, error) {\n\t// TODO: Not yet implemented.\n\treturn nil, status.Error(codes.Unimplemented, \"TODO: You should implement it!\")\n}\n{{end -}}\n"

// Service returns go-assets FileSystem
var Service = assets.NewFileSystem(map[string][]string{"/api/protos": []string{"{{.Path}}.proto.tmpl"}, "/app": []string{}, "/app/server": []string{"{{.Path}}_server.go.tmpl"}, "/": []string{}, "/api": []string{}}, map[string]*assets.File{
	"/app": &assets.File{
		Path:     "/app",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/app/server": &assets.File{
		Path:     "/app/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521908197, 1521908197000000000),
		Data:     nil,
	}, "/app/server/{{.Path}}_server.go.tmpl": &assets.File{
		Path:     "/app/server/{{.Path}}_server.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1521908197, 1521908197000000000),
		Data:     []byte(_Serviceef91c225ded973f86ca6b58050abcc766e9d41c3),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/api": &assets.File{
		Path:     "/api",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/api/protos": &assets.File{
		Path:     "/api/protos",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521908053, 1521908053000000000),
		Data:     nil,
	}, "/api/protos/{{.Path}}.proto.tmpl": &assets.File{
		Path:     "/api/protos/{{.Path}}.proto.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1521908053, 1521908053000000000),
		Data:     []byte(_Serviceba3c18adf9d77d5cbf8a40461aa014bc49f81edc),
	}}, "")
