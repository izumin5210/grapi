package generator

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
	"github.com/serenize/snaker"
	"github.com/spf13/afero"

	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module"
	"github.com/izumin5210/grapi/pkg/grapicmd/internal/module/generator/template"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

type serviceGenerator struct {
	baseGenerator
	rootDir string
}

func newServiceGenerator(fs afero.Fs, ui module.UI, rootDir string) module.ServiceGenerator {
	return &serviceGenerator{
		baseGenerator: newBaseGenerator(template.Service, fs, ui),
		rootDir:       rootDir,
	}
}

func (g *serviceGenerator) GenerateService(name string, methods ...string) error {
	data, err := g.createParams(name, methods)
	if err != nil {
		return errors.WithStack(err)
	}
	return g.Generate(g.rootDir, data)
}

func (g *serviceGenerator) DestroyService(name string) error {
	data, err := g.createParams(name, []string{})
	if err != nil {
		return errors.WithStack(err)
	}
	return g.Destroy(g.rootDir, data)
}

type nameParams struct {
	pluralCamel        string
	pluralCamelLower   string
	pluralSnake        string
	singularCamel      string
	singularCamelLower string
	singularSnake      string
}

type serviceParams struct {
	Path        string
	ServiceName string
	Methods     []serviceMethodParams
	Proto       serviceProtoParams
	PbGo        servicePbGoParams
	Go          serviceGoParams
}

type serviceProtoParams struct {
	Package  string
	Imports  []string
	Messages []serviceMethodMessage
}

type servicePbGoParams struct {
	PackageName string
	PackagePath string
}

type serviceGoParams struct {
	Package    string
	Imports    []string
	ServerName string
	StructName string
}

type serviceMethodsParams struct {
	Methods      []serviceMethodParams
	ProtoImports []string
	GoImports    []string
	Messages     []serviceMethodMessage
}

type serviceMethodParams struct {
	Method         string
	HTTP           serviceMethodHTTPParams
	requestCommon  string
	requestGo      string
	requestProto   string
	responseCommon string
	responseGo     string
	responseProto  string
}

func (p *serviceMethodParams) RequestGo(pkg string) string {
	if p.requestGo == "" {
		return pkg + "." + p.requestCommon
	}
	return p.requestGo
}

func (p *serviceMethodParams) RequestProto() string {
	if p.requestProto == "" {
		return p.requestCommon
	}
	return p.requestProto
}

func (p *serviceMethodParams) ResponseGo(pkg string) string {
	if p.responseGo == "" {
		return pkg + "." + p.responseCommon
	}
	return p.responseGo
}

func (p *serviceMethodParams) ResponseProto() string {
	if p.responseProto == "" {
		return p.responseCommon
	}
	return p.responseProto
}

type serviceMethodMessage struct {
	Name   string
	Fields []serviceMethodMessageField
}

type serviceMethodMessageField struct {
	Name     string
	Type     string
	Repeated bool
	Tag      uint
}

type serviceMethodHTTPParams struct {
	Method string
	Path   string
	Body   string
}

func (g *serviceGenerator) createParams(path string, methodNames []string) (*serviceParams, error) {
	// github.com/foo/bar
	importPath, err := fs.GetImportPath(g.rootDir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// path => baz/qux/quux
	path = strings.Replace(path, "-", "_", -1)

	// quux
	name := filepath.Base(path)

	nameParams := nameParams{
		pluralCamel:   inflection.Plural(snaker.SnakeToCamel(name)),
		singularCamel: inflection.Singular(snaker.SnakeToCamel(name)),
	}
	nameParams.pluralCamelLower = strings.ToLower(string(nameParams.pluralCamel[0])) + nameParams.pluralCamel[1:]
	nameParams.pluralSnake = snaker.CamelToSnake(nameParams.pluralCamel)
	nameParams.singularCamelLower = strings.ToLower(string(nameParams.singularCamel[0])) + nameParams.singularCamel[1:]
	nameParams.singularSnake = snaker.CamelToSnake(nameParams.singularCamel)

	// Quux
	serviceName := nameParams.singularCamel
	// quux
	localServiceName := strings.ToLower(string(serviceName[0])) + serviceName[1:]

	// baz/qux
	packagePath := filepath.Dir(path)
	// qux
	packageName := filepath.Base(packagePath)

	// api/baz/qux
	pbgoPackagePath := filepath.Join("api", packagePath)
	// qux_pb
	pbgoPackageName := filepath.Base(pbgoPackagePath) + "_pb"

	if packagePath == "." {
		packagePath = "server"
		packageName = packagePath
		pbgoPackagePath = "api"
		pbgoPackageName = pbgoPackagePath + "_pb"
	}

	protoPackageChunks := []string{}
	for _, pkg := range strings.Split(filepath.Join(importPath, "api", filepath.Dir(path)), "/") {
		chunks := strings.Split(strings.Replace(pkg, "-", "_", -1), ".")
		for i := len(chunks) - 1; i >= 0; i-- {
			protoPackageChunks = append(protoPackageChunks, chunks[i])
		}
	}
	// com.github.foo.bar.baz.qux
	protoPackage := strings.Join(protoPackageChunks, ".")

	protoImports := []string{
		"google/api/annotations.proto",
	}
	goImports := []string{
		"google.golang.org/grpc",
		"google.golang.org/grpc/codes",
		"google.golang.org/grpc/status",
	}

	methods := g.createMethodParams(nameParams, methodNames)

	protoImports = append(protoImports, methods.ProtoImports...)
	sort.Strings(protoImports)
	goImports = append(goImports, methods.GoImports...)
	sort.Strings(goImports)

	params := &serviceParams{
		Path:        path,
		ServiceName: serviceName,
		Methods:     methods.Methods,
		Proto: serviceProtoParams{
			Package:  protoPackage,
			Imports:  protoImports,
			Messages: methods.Messages,
		},
		PbGo: servicePbGoParams{
			PackageName: pbgoPackageName,
			PackagePath: filepath.Join(importPath, pbgoPackagePath),
		},
		Go: serviceGoParams{
			Package:    packageName,
			Imports:    goImports,
			ServerName: serviceName + "Service" + "Server",
			StructName: localServiceName + "Service" + "Server" + "Impl",
		},
	}

	return params, nil
}

func (g *serviceGenerator) createMethodParams(name nameParams, methods []string) (
	params serviceMethodsParams,
) {
	id := name.singularSnake + "_id"
	resource := &serviceMethodMessage{
		Name:   name.singularCamel,
		Fields: []serviceMethodMessageField{{Name: id, Type: "string", Tag: 1}},
	}

	basicMethods := [5]*serviceMethodParams{}
	customMethods := []serviceMethodParams{}
	basicMessages := [7]*serviceMethodMessage{}
	customMessages := []serviceMethodMessage{}

	for _, meth := range methods {
		switch strings.ToLower(meth) {
		case "list":
			methodName := "List" + name.pluralCamel
			reqName := methodName + "Request"
			respName := methodName + "Response"
			basicMethods[0] = &serviceMethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: respName,
				HTTP:           serviceMethodHTTPParams{Method: "get", Path: name.pluralSnake},
			}
			basicMessages[0] = resource
			basicMessages[1] = &serviceMethodMessage{Name: reqName}
			basicMessages[2] = &serviceMethodMessage{
				Name:   respName,
				Fields: []serviceMethodMessageField{{Name: name.pluralSnake, Type: name.singularCamel, Repeated: true, Tag: 1}},
			}
		case "get":
			methodName := "Get" + name.singularCamel
			reqName := methodName + "Request"
			basicMethods[1] = &serviceMethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: resource.Name,
				HTTP:           serviceMethodHTTPParams{Method: "get", Path: name.pluralSnake + "/{" + id + "}"},
			}
			basicMessages[0] = resource
			basicMessages[3] = &serviceMethodMessage{
				Name:   reqName,
				Fields: []serviceMethodMessageField{{Name: id, Type: "string", Tag: 1}},
			}
		case "create":
			methodName := "Create" + name.singularCamel
			reqName := methodName + "Request"
			basicMethods[2] = &serviceMethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: resource.Name,
				HTTP:           serviceMethodHTTPParams{Method: "post", Path: name.pluralSnake, Body: name.singularSnake},
			}
			basicMessages[0] = resource
			basicMessages[4] = &serviceMethodMessage{
				Name:   reqName,
				Fields: []serviceMethodMessageField{{Name: name.singularSnake, Type: name.singularCamel, Tag: 1}},
			}
		case "update":
			methodName := "Update" + name.singularCamel
			reqName := methodName + "Request"
			basicMethods[3] = &serviceMethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: resource.Name,
				HTTP:           serviceMethodHTTPParams{Method: "patch", Path: name.pluralSnake + "/{" + name.singularSnake + "." + id + "}", Body: name.singularSnake},
			}
			basicMessages[0] = resource
			basicMessages[5] = &serviceMethodMessage{
				Name:   reqName,
				Fields: []serviceMethodMessageField{{Name: name.singularSnake, Type: name.singularCamel, Tag: 1}},
			}
		case "delete":
			methodName := "Delete" + name.singularCamel
			reqName := methodName + "Request"
			basicMethods[4] = &serviceMethodParams{
				Method:        methodName,
				requestCommon: reqName,
				responseProto: "google.protobuf.Empty",
				responseGo:    "empty.Empty",
				HTTP:          serviceMethodHTTPParams{Method: "patch", Path: name.pluralSnake + "/{" + id + "}"},
			}
			basicMessages[6] = &serviceMethodMessage{
				Name:   reqName,
				Fields: []serviceMethodMessageField{{Name: id, Type: "string", Tag: 1}},
			}
			params.ProtoImports = append(params.ProtoImports, "google/protobuf/empty.proto")
			params.GoImports = append(params.GoImports, "github.com/golang/protobuf/ptypes/empty")
		default:
			methodName := snaker.SnakeToCamel(meth)
			reqName := methodName + "Request"
			respName := methodName + "Response"
			customMethods = append(customMethods, serviceMethodParams{
				Method:         methodName,
				requestCommon:  reqName,
				responseCommon: respName,
				HTTP:           serviceMethodHTTPParams{Method: "get", Path: name.pluralSnake + "/" + snaker.CamelToSnake(meth)},
			})
			customMessages = append(
				customMessages,
				serviceMethodMessage{Name: reqName},
				serviceMethodMessage{Name: respName},
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
