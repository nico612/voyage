package options

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/nico612/voyage/internal/pkg/server"
	"github.com/spf13/pflag"
	"time"
)

// JwtOptions JWT 配置
type JwtOptions struct {
	Realm      string        `json:"realm"       mapstructure:"realm"`       // jwt 标识
	Key        string        `json:"key"         mapstructure:"key"`         // 密钥
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout"`     // 过期时间
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh"` // 刷新时间
}

func NewJwtOptions() *JwtOptions {

	// 获取默认配置
	defaults := server.NewConfig()

	return &JwtOptions{
		Realm:      defaults.Jwt.Realm,
		Key:        defaults.Jwt.Key,
		Timeout:    defaults.Jwt.Timeout,
		MaxRefresh: defaults.Jwt.MaxRefresh,
	}
}

func (s *JwtOptions) ApplyTo(c *server.Config) error {
	c.Jwt = &server.JwtInfo{
		Realm:      s.Realm,
		Key:        s.Key,
		Timeout:    s.Timeout,
		MaxRefresh: s.MaxRefresh,
	}

	return nil
}

// Validate 验证配置
func (s *JwtOptions) Validate() []error {
	var errs []error

	if !govalidator.StringLength(s.Key, "6", "32") {
		errs = append(errs, fmt.Errorf("--secret-key must larger than 5 and little than 33"))
	}

	return errs
}

// AddFlags 将标志位添加到指定的 pflag.FlagSet
func (s *JwtOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.StringVar(&s.Realm, "jwt.realm", s.Realm, "Realm name to display to the user.")
	fs.StringVar(&s.Key, "jwt.key", s.Key, "Private key used to sign jwt token.")
	fs.DurationVar(&s.Timeout, "jwt.timeout", s.Timeout, "JWT token timeout.")

	fs.DurationVar(&s.MaxRefresh, "jwt.max-refresh", s.MaxRefresh, ""+
		"This field allows clients to refresh their token until MaxRefresh has passed.")
}
