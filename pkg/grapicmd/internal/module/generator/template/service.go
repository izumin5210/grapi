package template

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Service571266c56494500f3fd0d44c295886b67f7792c6 = "syntax = \"proto3\";\noption go_package = \"{{ .pbgoPackageName }}\";\npackage {{ .protoPackage }};\n\nimport \"google/api/annotations.proto\";\n\nservice {{ .serviceName }}Service {\t\n  rpc Get{{ .serviceName }} (Get{{ .serviceName }}Request) returns (Get{{ .serviceName }}Response) {\n    option (google.api.http) = {\n      get:  \"/{{ .path }}\"\n    };\n  }\n}\n\nmessage Get{{ .serviceName }}Request {\n}\n\nmessage Get{{ .serviceName }}Response {\n}\n"
var _Serviceade68dff2f92354600e62afbe061cceb4a0e52a6 = "package {{ .packageName }}\n\nimport (\n\t\"context\"\n\n\t\"google.golang.org/grpc\"\n\t\"google.golang.org/grpc/codes\"\n\t\"google.golang.org/grpc/status\"\n\n  {{ .pbgoPackageName }} \"{{ .importPath }}/{{ .pbgoPackagePath }}\"\n)\n\nvar (\n\t// Register{{ .serviceName }}ServiceHandler is a function to register card service handler to gRPC Gateway's mux.\n\tRegister{{ .serviceName }}ServiceHandler = {{ .pbgoPackageName }}.Register{{ .serviceName }}ServiceHandler\n)\n\n// Register{{ .serviceName }}ServiceServerFactory creates a function to register card service server impl to grpc.Server.\nfunc Register{{ .serviceName }}ServiceServerFactory() func(s *grpc.Server) {\n\treturn func(s *grpc.Server) {\n\t\t{{ .pbgoPackageName }}.Register{{ .serviceName }}ServiceServer(s, New())\n\t}\n}\n\n// New creates a new {{ .serviceName }}ServiceServer instance.\nfunc New() {{ .pbgoPackageName }}.{{ .serviceName }}ServiceServer {\n\treturn &{{ .localServiceName }}ServiceServerImpl{}\n}\n\ntype {{ .localServiceName }}ServiceServerImpl struct {\n}\n\nfunc (s *{{ .localServiceName }}ServiceServerImpl) Get{{ .serviceName }}(ctx context.Context, req *{{ .pbgoPackageName }}.Get{{ .serviceName }}Request) (*{{ .pbgoPackageName }}.Get{{ .serviceName }}Response, error) {\n\t// TODO: Not yet implemented.\n\treturn nil, status.Error(codes.Unimplemented, \"TODO: You should implement it!\")\n}\n"

// Service returns go-assets FileSystem
var Service = assets.NewFileSystem(map[string][]string{"/app": []string{}, "/app/server": []string{"{{ .path }}_server.go.tmpl"}, "/": []string{}, "/api": []string{}, "/api/protos": []string{"{{ .path }}.proto.tmpl"}}, map[string]*assets.File{
	"/app/server/{{ .path }}_server.go.tmpl": &assets.File{
		Path:     "/app/server/{{ .path }}_server.go.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1521888919, 1521888919000000000),
		Data:     []byte(_Serviceade68dff2f92354600e62afbe061cceb4a0e52a6),
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
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/api/protos/{{ .path }}.proto.tmpl": &assets.File{
		Path:     "/api/protos/{{ .path }}.proto.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     []byte(_Service571266c56494500f3fd0d44c295886b67f7792c6),
	}, "/app": &assets.File{
		Path:     "/app",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1520753819, 1520753819000000000),
		Data:     nil,
	}, "/app/server": &assets.File{
		Path:     "/app/server",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1521888919, 1521888919000000000),
		Data:     nil,
	}}, "")
