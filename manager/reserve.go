package manager

import "context"

type ReserveConfig struct {
	BaseHost       string `yaml:"base_host" json:"base_host,omitempty" arg:"" name:"base_host" help:"main service host"`
	BasePort       int    `yaml:"base_port" json:"base_port,omitempty" arg:"" name:"base_port" help:"main service port"`
	ReserveHost    string `yaml:"reserve_host" json:"reserve_host,omitempty" arg:"" name:"reserve_host" help:"reserve service host"`
	ReservePort    int    `yaml:"reserve_port" json:"reserve_port,omitempty" arg:"" name:"reserve_port" help:"reserve service port"`
	ListenPort     int    `yaml:"listen_port" json:"listen_port,omitempty" arg:"" name:"listen_port" help:"port to listen" default:"9000"`
	ConnectTimeout int    `yaml:"timeout" json:"timeout,omitempty" arg:"" name:"connect_timeout" help:"connect_timeout" default:"3"`
	MaxIdleSeconds int    `yaml:"max_idle_seconds" json:"max_idle_seconds,omitempty" arg:"" name:"max_idle_seconds" help:"max_idle_seconds" default:"600"`
}

type Reserve interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Error() error
}
