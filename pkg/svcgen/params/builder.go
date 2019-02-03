package params

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/serenize/snaker"

	"github.com/izumin5210/grapi/pkg/cli"
	gencmdutil "github.com/izumin5210/grapi/pkg/gencmd/util"
)

type Builder interface {
	Build(path string, resName string, methodNames []string) (*Params, error)
}

func NewBuilder(rootDir cli.RootDir, protoDir, protoOutDir, serverDir string, pkgName string) Builder {
	if protoDir == "" {
		protoDir = filepath.Join("api", "protos")
	}
	if protoOutDir == "" {
		protoOutDir = filepath.Join("api")
	}
	if serverDir == "" {
		serverDir = filepath.Join("app", "server")
	}
	return &builderImpl{
		rootDir:     rootDir,
		protoDir:    protoDir,
		protoOutDir: protoOutDir,
		serverDir:   serverDir,
		pkgName:     pkgName,
	}
}

type builderImpl struct {
	rootDir     cli.RootDir
	protoDir    string
	protoOutDir string
	serverDir   string
	pkgName     string
}

func (b *builderImpl) Build(path string, resName string, methodNames []string) (*Params, error) {
	protoParams, err := gencmdutil.BuildProtoParams(path, b.rootDir, b.protoOutDir, b.pkgName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// path => baz/qux/quux
	path = protoParams.Proto.Path

	// quux
	name := filepath.Base(path)

	names := gencmdutil.Inflect(name)

	// Quux
	serviceName := names.Camel.Singular
	// quux
	localServiceName := strings.ToLower(string(serviceName[0])) + serviceName[1:]

	// baz/qux
	packagePath := filepath.Dir(path)
	// qux
	packageName := filepath.Base(packagePath)

	if packagePath == "." {
		packagePath = filepath.Base(b.serverDir)
		packageName = packagePath
	}

	protoImports := []string{
		"google/api/annotations.proto",
	}
	goImports := []string{
		"github.com/izumin5210/grapi/pkg/grapiserver",
		"google.golang.org/grpc/codes",
		"google.golang.org/grpc/status",
	}

	resNames := names
	if resName != "" {
		resNames = gencmdutil.Inflect(resName)
	}
	methods := b.buildMethodParams(resNames, methodNames)

	protoImports = append(protoImports, methods.ProtoImports...)
	sort.Strings(protoImports)
	goImports = append(goImports, methods.GoImports...)
	sort.Strings(goImports)

	params := &Params{
		ProtoDir:    b.protoDir,
		ProtoOutDir: b.protoOutDir,
		ServerDir:   b.serverDir,
		Path:        path,
		ServiceName: serviceName,
		Methods:     methods.Methods,
		Proto: ProtoParams{
			Package:  protoParams.Proto.Package,
			Imports:  protoImports,
			Messages: methods.Messages,
		},
		PbGo: PbGoParams{
			PackagePath: protoParams.PbGo.Package,
			PackageName: protoParams.PbGo.ImportName,
		},
		Go: GoParams{
			Package:    packageName,
			Imports:    goImports,
			ServerName: serviceName + "Service" + "Server",
			StructName: localServiceName + "Service" + "Server" + "Impl",
		},
	}

	return params, nil
}

func (b *builderImpl) buildMethodParams(name gencmdutil.String, methods []string) (
	params MethodsParams,
) {
	id := name.Snake.Singular + "_id"
	resource := &MethodMessage{
		Name:   name.Camel.Singular,
		Fields: []MethodMessageField{{Name: id, Type: "string", Tag: 1}},
	}

	basicMethods := [5]*MethodParams{}
	customMethods := []MethodParams{}
	basicMessages := [7]*MethodMessage{}
	customMessages := []MethodMessage{}

	for _, meth := range methods {
		switch strings.ToLower(meth) {
		case "list":
			methodName := "List" + name.Camel.Plural
			reqName := methodName + "Request"
			respName := methodName + "Response"
			basicMethods[0] = &MethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: respName,
				HTTP:           MethodHTTPParams{Method: "get", Path: name.Snake.Plural},
			}
			basicMessages[0] = resource
			basicMessages[1] = &MethodMessage{Name: reqName}
			basicMessages[2] = &MethodMessage{
				Name:   respName,
				Fields: []MethodMessageField{{Name: name.Snake.Plural, Type: name.Camel.Singular, Repeated: true, Tag: 1}},
			}
		case "get":
			methodName := "Get" + name.Camel.Singular
			reqName := methodName + "Request"
			basicMethods[1] = &MethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: resource.Name,
				HTTP:           MethodHTTPParams{Method: "get", Path: name.Snake.Plural + "/{" + id + "}"},
			}
			basicMessages[0] = resource
			basicMessages[3] = &MethodMessage{
				Name:   reqName,
				Fields: []MethodMessageField{{Name: id, Type: "string", Tag: 1}},
			}
		case "create":
			methodName := "Create" + name.Camel.Singular
			reqName := methodName + "Request"
			basicMethods[2] = &MethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: resource.Name,
				HTTP:           MethodHTTPParams{Method: "post", Path: name.Snake.Plural, Body: name.Snake.Singular},
			}
			basicMessages[0] = resource
			basicMessages[4] = &MethodMessage{
				Name:   reqName,
				Fields: []MethodMessageField{{Name: name.Snake.Singular, Type: name.Camel.Singular, Tag: 1}},
			}
		case "update":
			methodName := "Update" + name.Camel.Singular
			reqName := methodName + "Request"
			basicMethods[3] = &MethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: resource.Name,
				HTTP:           MethodHTTPParams{Method: "patch", Path: name.Snake.Plural + "/{" + name.Snake.Singular + "." + id + "}", Body: name.Snake.Singular},
			}
			basicMessages[0] = resource
			basicMessages[5] = &MethodMessage{
				Name:   reqName,
				Fields: []MethodMessageField{{Name: name.Snake.Singular, Type: name.Camel.Singular, Tag: 1}},
			}
		case "delete":
			methodName := "Delete" + name.Camel.Singular
			reqName := methodName + "Request"
			basicMethods[4] = &MethodParams{
				Method:        methodName,
				requestCommon: reqName,
				responseProto: "google.protobuf.Empty",
				responseGo:    "empty.Empty",
				HTTP:          MethodHTTPParams{Method: "delete", Path: name.Snake.Plural + "/{" + id + "}"},
			}
			basicMessages[6] = &MethodMessage{
				Name:   reqName,
				Fields: []MethodMessageField{{Name: id, Type: "string", Tag: 1}},
			}
			params.ProtoImports = append(params.ProtoImports, "google/protobuf/empty.proto")
			params.GoImports = append(params.GoImports, "github.com/golang/protobuf/ptypes/empty")
		default:
			methodName := snaker.SnakeToCamel(meth)
			reqName := methodName + "Request"
			respName := methodName + "Response"
			customMethods = append(customMethods, MethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: respName,
				HTTP:           MethodHTTPParams{Method: "get", Path: name.Snake.Plural + "/" + snaker.CamelToSnake(meth)},
			})
			customMessages = append(
				customMessages,
				MethodMessage{Name: reqName},
				MethodMessage{Name: respName},
			)
		}
	}

	for _, meth := range basicMethods {
		if meth != nil {
			params.Methods = append(params.Methods, *meth)
		}
	}
	for _, msg := range basicMessages {
		if msg != nil {
			params.Messages = append(params.Messages, *msg)
		}
	}
	for _, meth := range customMethods {
		params.Methods = append(params.Methods, meth)
	}
	for _, msg := range customMessages {
		params.Messages = append(params.Messages, msg)
	}
	return
}
