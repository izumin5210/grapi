package params

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
	"github.com/pkg/errors"
	"github.com/serenize/snaker"
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
	// github.com/foo/bar
	importPath, err := fs.GetImportPath(b.rootDir.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// path => baz/qux/quux
	path = strings.Replace(path, "-", "_", -1)

	// quux
	name := filepath.Base(path)

	names := inflect(name)

	// Quux
	serviceName := names.singularCamel
	// quux
	localServiceName := strings.ToLower(string(serviceName[0])) + serviceName[1:]

	// baz/qux
	packagePath := filepath.Dir(path)
	// qux
	packageName := filepath.Base(packagePath)

	// api/baz/qux
	pbgoPackagePath := filepath.Join(b.protoOutDir, packagePath)
	// qux_pb
	pbgoPackageName := filepath.Base(pbgoPackagePath) + "_pb"

	if packagePath == "." {
		packagePath = filepath.Base(b.serverDir)
		packageName = packagePath
		pbgoPackagePath = b.protoOutDir
		pbgoPackageName = filepath.Base(pbgoPackagePath) + "_pb"
	}

	protoPackage := b.pkgName
	if protoPackage == "" {
		protoPackageChunks := []string{}
		for _, pkg := range strings.Split(filepath.ToSlash(filepath.Join(importPath, b.protoOutDir)), "/") {
			chunks := strings.Split(pkg, ".")
			for i := len(chunks) - 1; i >= 0; i-- {
				protoPackageChunks = append(protoPackageChunks, chunks[i])
			}
		}
		// com.github.foo.bar.baz.qux
		protoPackage = strings.Join(protoPackageChunks, ".")
	}
	if dir := filepath.Dir(path); dir != "." {
		protoPackage = protoPackage + "." + strings.Replace(dir, string(filepath.Separator), ".", -1)
	}
	protoPackage = strings.Replace(protoPackage, "-", "_", -1)

	protoImports := []string{
		"google/api/annotations.proto",
	}
	goImports := []string{
		"github.com/izumin5210/grapi/pkg/grapiserver",
		"google.golang.org/grpc/codes",
		"google.golang.org/grpc/status",
	}
	goTestImports := []string{}

	resNames := names
	if resName != "" {
		resNames = inflect(resName)
	}
	methods := b.buildMethodParams(resNames, methodNames)

	protoImports = append(protoImports, methods.ProtoImports...)
	sort.Strings(protoImports)
	goImports = append(goImports, methods.GoImports...)
	sort.Strings(goImports)
	goTestImports = append(goTestImports, methods.GoImports...)
	sort.Strings(goTestImports)

	params := &Params{
		ProtoDir:    b.protoDir,
		ProtoOutDir: b.protoOutDir,
		ServerDir:   b.serverDir,
		Path:        path,
		ServiceName: serviceName,
		Methods:     methods.Methods,
		Proto: ProtoParams{
			Package:  protoPackage,
			Imports:  protoImports,
			Messages: methods.Messages,
		},
		PbGo: PbGoParams{
			PackageName: pbgoPackageName,
			PackagePath: filepath.ToSlash(filepath.Join(importPath, pbgoPackagePath)),
		},
		Go: GoParams{
			Package:     packageName,
			Imports:     goImports,
			TestImports: goTestImports,
			ServerName:  serviceName + "Service" + "Server",
			StructName:  localServiceName + "Service" + "Server" + "Impl",
		},
	}

	return params, nil
}

func (b *builderImpl) buildMethodParams(name inflectableString, methods []string) (
	params MethodsParams,
) {
	id := name.singularSnake + "_id"
	resource := &MethodMessage{
		Name:   name.singularCamel,
		Fields: []MethodMessageField{{Name: id, Type: "string", Tag: 1}},
	}

	basicMethods := [5]*MethodParams{}
	customMethods := []MethodParams{}
	basicMessages := [7]*MethodMessage{}
	customMessages := []MethodMessage{}

	for _, meth := range methods {
		switch strings.ToLower(meth) {
		case "list":
			methodName := "List" + name.pluralCamel
			reqName := methodName + "Request"
			respName := methodName + "Response"
			basicMethods[0] = &MethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: respName,
				HTTP:           MethodHTTPParams{Method: "get", Path: name.pluralSnake},
			}
			basicMessages[0] = resource
			basicMessages[1] = &MethodMessage{Name: reqName}
			basicMessages[2] = &MethodMessage{
				Name:   respName,
				Fields: []MethodMessageField{{Name: name.pluralSnake, Type: name.singularCamel, Repeated: true, Tag: 1}},
			}
		case "get":
			methodName := "Get" + name.singularCamel
			reqName := methodName + "Request"
			basicMethods[1] = &MethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: resource.Name,
				HTTP:           MethodHTTPParams{Method: "get", Path: name.pluralSnake + "/{" + id + "}"},
			}
			basicMessages[0] = resource
			basicMessages[3] = &MethodMessage{
				Name:   reqName,
				Fields: []MethodMessageField{{Name: id, Type: "string", Tag: 1}},
			}
		case "create":
			methodName := "Create" + name.singularCamel
			reqName := methodName + "Request"
			basicMethods[2] = &MethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: resource.Name,
				HTTP:           MethodHTTPParams{Method: "post", Path: name.pluralSnake, Body: name.singularSnake},
			}
			basicMessages[0] = resource
			basicMessages[4] = &MethodMessage{
				Name:   reqName,
				Fields: []MethodMessageField{{Name: name.singularSnake, Type: name.singularCamel, Tag: 1}},
			}
		case "update":
			methodName := "Update" + name.singularCamel
			reqName := methodName + "Request"
			basicMethods[3] = &MethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: resource.Name,
				HTTP:           MethodHTTPParams{Method: "patch", Path: name.pluralSnake + "/{" + name.singularSnake + "." + id + "}", Body: name.singularSnake},
			}
			basicMessages[0] = resource
			basicMessages[5] = &MethodMessage{
				Name:   reqName,
				Fields: []MethodMessageField{{Name: name.singularSnake, Type: name.singularCamel, Tag: 1}},
			}
		case "delete":
			methodName := "Delete" + name.singularCamel
			reqName := methodName + "Request"
			basicMethods[4] = &MethodParams{
				Method:        methodName,
				requestCommon: reqName,
				responseProto: "google.protobuf.Empty",
				responseGo:    "empty.Empty",
				HTTP:          MethodHTTPParams{Method: "delete", Path: name.pluralSnake + "/{" + id + "}"},
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
				HTTP:           MethodHTTPParams{Method: "get", Path: name.pluralSnake + "/" + snaker.CamelToSnake(meth)},
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
