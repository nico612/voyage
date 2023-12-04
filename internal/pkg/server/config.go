package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/pkg/log"
	"github.com/nico612/voyage/pkg/utils/homedir"
	"github.com/spf13/viper"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// RecommendedHomeDir defines the default directory used to place all iam service configurations.
	RecommendedHomeDir = ".adminsrv"

	// RecommendedEnvPrefix defines the ENV prefix used by all iam service.
	RecommendedEnvPrefix = "VOYAGE"
)

// Config 是用于配置 GenericAPIServer 的结构。其成员大致按照对组件的重要性排序。
type Config struct {
	SecureServing   *SecureServingInfo   // https 配置
	InsecureServing *InsecureServingInfo // http 配置
	Jwt             *JwtInfo             // jwt 配置
	Mode            string               // 运行模式
	Middlewares     []string             // 中间件
	Healthz         bool                 // 健康检查
	EnableProfiling bool                 // 性能监控
	EnableMetrics   bool                 // 指标监控
}

// CertKey 包含与证书相关的配置项。
type CertKey struct {
	// CertFile 是一个包含PEM编码证书的文件，可能还包括完整的证书链。
	CertFile string
	// KeyFile 是一个包含PEM编码私钥的文件，对应于CertFile指定的证书。
	KeyFile string
}

// SecureServingInfo 保存TLS服务器的配置信息。
type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

// Address 将主机IP地址和主机端口号拼接成一个地址字符串，例如：0.0.0.0:8443。
func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

// InsecureServingInfo 保存不安全的 HTTP 服务器的配置信息。
type InsecureServingInfo struct {
	Address string
}

// JwtInfo 定义了用于创建JWT身份验证中间件的JWT字段。
type JwtInfo struct {
	// defaults to "adminsrv jwt"
	Realm string
	// defaults to empty
	Key string
	// defaults to one hour
	Timeout time.Duration
	// defaults to zero
	MaxRefresh time.Duration
}

// NewConfig 默认配置
func NewConfig() *Config {
	return &Config{
		Healthz:         true,
		Mode:            gin.ReleaseMode,
		Middlewares:     []string{},
		EnableProfiling: true,
		EnableMetrics:   true,
		Jwt: &JwtInfo{
			Realm:      "adminsrv jwt",
			Timeout:    1 * time.Hour,
			MaxRefresh: 1 * time.Hour,
		},
	}
}

// Complete 定义了一个 Complete 方法，用于确保配置信息的完整性。它返回一个 CompletedConfig 结构，该结构包含已完成的配置，并对配置信息进行必要的补全。
// 这种方法设计的目的是确保配置信息完整且满足所有必要的条件，以便在使用它们之前可以信任它们。
func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

// CompletedConfig  GenericAPIServer 已完成配置. 这里这样设计主要是为了方便扩展
type CompletedConfig struct {
	*Config
}

// New 根据完整配置返回一个新的 GenericAPIServer 实例
func (c CompletedConfig) New() (*GenericAPIServer, error) {
	// setMode before gin.New()
	gin.SetMode(c.Mode)

	s := &GenericAPIServer{
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		healthz:             c.Healthz,
		enableMetrics:       c.EnableMetrics,
		enableProfiling:     c.EnableProfiling,
		middlewares:         c.Middlewares,
		Engine:              gin.New(),
	}

	initGenericAPIServer(s)

	return s, nil
}

// LoadConfig 读取配置文件和已设置的环境变量。
func LoadConfig(cfg string, defaultName string) {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(homedir.HomeDir(), RecommendedHomeDir))
		viper.AddConfigPath("/etc/adminsrv")
		viper.SetConfigName(defaultName)
	}

	// Use config file from the flag.
	viper.SetConfigType("yaml")              // set the type of the configuration to yaml.
	viper.AutomaticEnv()                     // read in environment variables that match.
	viper.SetEnvPrefix(RecommendedEnvPrefix) // set ENVIRONMENT variables prefix to IAM.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("WARNING: viper failed to discover and load the configuration file: %s", err.Error())
	}
}
