package options

import (
	"github.com/nico612/voyage/internal/pkg/server"
	"github.com/spf13/pflag"
)

const EnvModel = "ENV_MODEL"

const (
	DebugModel   = "debug"
	ReleaseModel = "release"
	TestModel    = "test"
)

// ServerRunOptions 服务通用配置
type ServerRunOptions struct {
	Mode          string   `json:"mode"        mapstructure:"mode"`
	Healthz       bool     `json:"healthz"     mapstructure:"healthz"`
	Middlewares   []string `json:"middlewares" mapstructure:"middlewares"`            // 中间件，如果需要根据配置来配置中间件
	UseMultipoint bool     `json:"use-multipoint" mapstructure:"use-multipoint"`      // 多点登录
	IplimitCount  int32    `json:"iplimit-count"        mapstructure:"iplimit-count"` // ip 限制数量
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		Mode:          ReleaseModel,
		Healthz:       true,
		UseMultipoint: false,
		IplimitCount:  15000,
	}
}

// ApplyTo 应用到 api server 配置中
func (s *ServerRunOptions) ApplyTo(c *server.Config) error {
	c.Mode = s.Mode
	c.Healthz = s.Healthz
	c.Middlewares = s.Middlewares

	return nil
}

// Validate 检查 ServerRunOptions.
func (s *ServerRunOptions) Validate() []error {
	errors := []error{}

	return errors
}

// AddFlags 向 APIServer 添加 标志位到 指定的 FlagSet 中
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs.StringVar(&s.Mode, "server.mode", s.Mode, ""+
		"Start the server in a specified server mode. Supported server mode: debug, test, release.")

	fs.BoolVar(&s.Healthz, "server.healthz", s.Healthz, ""+
		"Add self readiness check and install /healthz router.")

	fs.StringSliceVar(&s.Middlewares, "server.middlewares", s.Middlewares, ""+
		"List of allowed middlewares for server, comma separated. If this list is empty default middlewares will be used.")
}
