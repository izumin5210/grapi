package grapicmd

type ProtocConfig struct {
	ImportDirs []string `mapstructure:"import_dirs"`
	ProtosDir  string   `mapstructure:"protos_dir"`
	OutDir     string   `mapstructure:"out_dir"`
	Plugins    []*ProtocPlugin
}

// ProtocPlugin contains args and plugin name for using in protoc command.
type ProtocPlugin struct {
	Path string
	Name string
	Args map[string]interface{}
}
