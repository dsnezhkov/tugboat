package defs

type Config struct {
	GlobalOpts     GlobalConfig    `hcl:"global,block"`
	VectorBlock    VectorConfig    `hcl:"vectors,block"`
	ComponentBlock ComponentConfig `hcl:"components,block"`
	FlowLinkBlock  FlowLinkConfig  `hcl:"flows,block"`
}
type GlobalConfig struct {
	LoggerOpts  LoggerConfig    `hcl:"logger,block"`
	PayloadOpts []PayloadConfig `hcl:"payloads,block"`
}
type LoggerConfig struct {
	Level    string `hcl:"level"`
	Location string `hcl:"location"`
}
type PayloadConfig struct {
	Variant  string `hcl:"variant,label"`
	Location string `hcl:"location"`
}

type VectorConfig struct {
	Vectors []Vector `hcl:"vector,block"`
}
type ComponentConfig struct {
	Components []Component `hcl:"component,block"`
}
type Vector struct {
	Name           string   `hcl:"name,label"`
	ComponentLinks []string `hcl:"components"`
}
type Component struct {
	Name          string   `hcl:"name,label"`
	Description   string   `hcl:"desc"`
	Type          string   `hcl:"type"`
	Active        bool     `hcl:"active"`
	Timeout       int      `hcl:"timeout,optional"`
	Directive     string   `hcl:"directive,optional"`
	DirectiveOpts []string `hcl:"directive_opts,optional"`
	Data          []string `hcl:"default_data,optional"`
	Modules       []Module `hcl:"module,block"`
}

type Module struct {
	Name  string `hcl:"module_name,label"`
	Loads []Load `hcl:"load,block"`
}
type Load struct {
	Name          string         `hcl:"name,label"`
	Identifier    string         `hcl:"identifier,optional"`
	DataType      string         `hcl:"data_type"`
	NativeFormat  string         `hcl:"native_format"`
	Encapsulation *Encapsulation `hcl:"encapsulation,block"` // * if block can be optional
	Source        string         `hcl:"source"`
}

type Encapsulation struct {
	StoredEncrypted  *StorageEncrypted  `hcl:"encrypted,block"`
	StoredCompressed *StorageCompressed `hcl:"compressed,block"`
	StoredEncoded    *StorageEncoded    `hcl:"encoded,block"`
	Order            string             `hcl:"encapsulation_order,optional"`
}

type StorageEncrypted struct {
	Algorithm string `hcl:"algorithm"`
	Key       string `hcl:"key"`
}
type StorageCompressed struct {
	Variant string `hcl:"variant"`
}
type StorageEncoded struct {
	Variant string `hcl:"variant"`
}
type FlowLinkConfig struct {
	Head      string     `hcl:"head"`
	FlowLinks []FlowLink `hcl:"from_vector,block"`
}
type FlowLink struct {
	Name        string   `hcl:"name,label"`
	VectorNames []string `hcl:"to_vectors"`
}
