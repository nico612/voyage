package config

import (
	"encoding/json"
	"github.com/nico612/voyage/pkg/utils/iduitl"

	"github.com/nico612/voyage/internal/pkg/options"
	cliflag "github.com/nico612/voyage/pkg/cli/flag"
	"github.com/nico612/voyage/pkg/log"
)

type Config struct {
	Server     *options.ServerRunOptions `json:"server" mapstructure:"server"`
	Grpc       *options.GRPCOptions      `json:"grpc" mapstructure:"grpc"`
	JwtOptions *options.JwtOptions       `json:"jwt"      mapstructure:"jwt"`
	Log        *log.Options              `json:"log"      mapstructure:"log"`
}

func NewConfig() *Config {
	return &Config{
		Server:     options.NewServerRunOptions(),
		Grpc:       options.NewGRPCOptions(),
		JwtOptions: options.NewJwtOptions(),
		Log:        log.NewOptions(),
	}
}

func (c *Config) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}

// Flags 可以在这里返回指定的命令行
func (o *Config) Flags() (fss cliflag.NamedFlagSets) {
	return fss
}

// Validate 检查所有配置并返回所有错误的集合
func (o *Config) Validate() []error {
	var errs []error

	errs = append(errs, o.JwtOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	// ... 检查其他配置

	return errs
}

// Complete set default Options.
func (c *Config) Complete() error {
	if c.JwtOptions.Key == "" {
		c.JwtOptions.Key = iduitl.NewSecretKey()
	}
	return nil
}
