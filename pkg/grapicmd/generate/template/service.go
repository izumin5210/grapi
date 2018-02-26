package template

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Service571266c56494500f3fd0d44c295886b67f7792c6 = "syntax = \"proto3\";\noption go_package = \"{{ .pbgoPackageName }}\";\npackage {{ .protoPackage }};\n\nimport \"google/api/annotations.proto\";\n\nservice {{ .serviceName }}Service {\t\n  rpc Get{{ .serviceName }} (Get{{ .serviceName }}Request) returns (Get{{ .serviceName }}Response) {\n    option (google.api.http) = {\n      get:  \"/{{ .path }}\"\n    };\n  }\n}\n\nmessage Get{{ .serviceName }}Request {\n}\n\nmessage Get{{ .serviceName }}Response {\n}\n"
var _Service0cea4966b52adb07bedd06cbe7ba96b91d5c98e2 = "package {{ .packageName }}\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\n\t\"google.golang.org/grpc\"\n\t\"google.golang.org/grpc/codes\"\n\t\"google.golang.org/grpc/status\"\n\n\t\"{{ .importPath }}/{{ .pbgoPackagePath }}\"\n)\n\n// Register{{ .serviceName }}ServiceServerFactory creates a function to register card service server impl to grpc.Server.\nfunc Register{{ .serviceName }}ServiceServerFactory() func(s *grpc.Server) {\n\treturn func(s *grpc.Server) {\n\t\t{{ .packageName }}.Register{{ .serviceName }}ServiceServer(s, New())\n\t}\n}\n\n// New creates a new {{ .serviceName }}ServiceServer instance.\nfunc New() {{ .pbgoPackageName }}.{{ .serviceName }}ServiceServer {\n\treturn &{{ .name }}ServiceServerImpl{}\n}\n\ntype {{ .name }}ServiceServerImpl struct {\n}\n\nfunc (s *{{ .name }}ServiceServerImpl) Get{{ .serviceName }}(ctx context.Context, req *{{ .pbgoPackageName }}.Get{{ .serviceName }}Request) (*{{ .pbgoPackageName }}.Get{{ .serviceName }}Response, error) {\n\treturn nil, status.Error(codes.Unimplemented, \"not yet implemented\")\n}\n"

// Service returns go-assets FileSystem
var Service = assets.NewFileSystem(map[string][]string{"/": []string{}, "/api": []string{}, "/api/protos": []string{"{{ .path }}.proto.tmpl"}, "/app": []string{}, "/app/server": []string{"{{ .path }}_service_server.go.tmpl"}}, map[string]*assets.File{
	"/app/server/{{ .path }}_service_server.go.tmpl": &assets.File{
		Path:     "/app/server/{{ .path }}_service_server.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519555491, 1519555491000000000),
		Data:     []byte(_Service0cea4966b52adb07bedd06cbe7ba96b91d5c98e2),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519541622, 1519541622000000000),
		Data:     nil,
	}, "/api": &assets.File{
		Path:     "/api",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519541622, 1519541622000000000),
		Data:     nil,
	}, "/api/protos": &assets.File{
		Path:     "/api/protos",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519554903, 1519554903000000000),
		Data:     nil,
	}, "/api/protos/{{ .path }}.proto.tmpl": &assets.File{
		Path:     "/api/protos/{{ .path }}.proto.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1519554903, 1519554903000000000),
		Data:     []byte(_Service571266c56494500f3fd0d44c295886b67f7792c6),
	}, "/app": &assets.File{
		Path:     "/app",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519541622, 1519541622000000000),
		Data:     nil,
	}, "/app/server": &assets.File{
		Path:     "/app/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1519555491, 1519555491000000000),
		Data:     nil,
	}}, "")
