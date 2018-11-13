package params

type Params struct {
	ProtoDir    string
	ProtoOutDir string
	ServerDir   string
	Path        string
	ServiceName string
	Methods     []MethodParams
	Proto       ProtoParams
	PbGo        PbGoParams
	Go          GoParams
}

type ProtoParams struct {
	Package  string
	Imports  []string
	Messages []MethodMessage
}

type PbGoParams struct {
	PackageName string
	PackagePath string
}

type GoParams struct {
	Package     string
	Imports     []string
	TestImports []string
	ServerName  string
	StructName  string
}

type MethodsParams struct {
	Methods      []MethodParams
	ProtoImports []string
	GoImports    []string
	Messages     []MethodMessage
}

type MethodParams struct {
	Method         string
	HTTP           MethodHTTPParams
	requestCommon  string
	requestGo      string
	requestProto   string
	responseCommon string
	responseGo     string
	responseProto  string
}

func (p *MethodParams) RequestGo(pkg string) string {
	if p.requestGo == "" {
		return pkg + "." + p.requestCommon
	}
	return p.requestGo
}

func (p *MethodParams) RequestProto() string {
	if p.requestProto == "" {
		return p.requestCommon
	}
	return p.requestProto
}

func (p *MethodParams) ResponseGo(pkg string) string {
	if p.responseGo == "" {
		return pkg + "." + p.responseCommon
	}
	return p.responseGo
}

func (p *MethodParams) ResponseProto() string {
	if p.responseProto == "" {
		return p.responseCommon
	}
	return p.responseProto
}

type MethodMessage struct {
	Name   string
	Fields []MethodMessageField
}

type MethodMessageField struct {
	Name     string
	Type     string
	Repeated bool
	Tag      uint
}

type MethodHTTPParams struct {
	Method string
	Path   string
	Body   string
}
