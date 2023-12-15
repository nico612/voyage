package config

import "github.com/nico612/voyage/internal/adminsrv/options"

var cfg *Config

// Config 服务运行配置文件
type Config struct {
	*options.Options
}

// CreateConfigFromOptions 根据配置选项创建配置
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {

	cfg = &Config{opts}
	return cfg, nil
}

func GetAppConfig() *Config {
	return cfg
}
