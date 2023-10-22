package manager

import "context"

type ForwardExtraHeader struct {
	Name  string `yaml:"name" json:"name,omitempty" arg:"" name:"name" help:"header name"`
	Value string `yaml:"value" json:"value,omitempty" arg:"" name:"value" help:"header value"`
}

type ForwardConfig struct {
	Target           string               `yaml:"target" json:"target,omitempty" arg:"" name:"target" help:"target"`
	ListenPort       int                  `yaml:"listen_port" json:"listen_port,omitempty" arg:"" name:"listen_port" help:"listen_port"`
	SourcePathPrefix string               `yaml:"source_path_prefix" json:"source_path_prefix,omitempty" arg:"" name:"source_path_prefix" help:"source path prefix"`
	Headers          []ForwardExtraHeader `yaml:"headers" json:"headers,omitempty"`
}

type Forwarder interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Error() error
}
