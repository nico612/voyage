package options

import (
	"github.com/spf13/pflag"
	"time"
)

// RedisOptions 定义 single Redis 配置
type RedisOptions struct {
	Addr         string        `json:"addr" mapstructure:"addr"`                   // Redis 连接地址
	Password     string        `json:"password" mapstructure:"password"`           // 密码
	DB           int           `json:"db" mapstructure:"db"`                       // Redis 数据库索引。
	ReadTimeout  time.Duration `json:"read-timeout" mapstructure:"read-timeout"`   // 读取超时时间
	WriteTimeout time.Duration `json:"write-timeout" mapstructure:"write-timeout"` // 写入超时时间
	DialTimeout  time.Duration `json:"dial-timeout" mapstructure:"dial-timeout"`   // 连接超时时间。
}

// NewRedisOptions create a `zero` value instance.
func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Addr:         "127.0.0.1",
		Password:     "",
		DB:           0,
		DialTimeout:  2 * time.Second,
		WriteTimeout: 300 * time.Millisecond,
		ReadTimeout:  300 * time.Millisecond,
	}
}

// Validate verifies flags passed to RedisOptions.
func (o *RedisOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds flags related to single redis for a specific APIServer to the specified FlagSet.
func (o *RedisOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Addr, "single.addr", o.Password, "Optional auth password for Redis db.")

	fs.StringVar(&o.Password, "single.password", o.Password, "Optional auth password for Redis db.")

	fs.IntVar(&o.DB, "single.database", o.DB, ""+
		"By default, the database is 0. Setting the database is not supported with single cluster. "+
		"As such, if you have --single.enable-cluster=true, then this value should be omitted or explicitly set to 0.")

	fs.DurationVar(&o.DialTimeout, "single.timeout", o.DialTimeout, "Timeout (in seconds) when connecting to single service.")
	fs.DurationVar(&o.ReadTimeout, "single.readTimeout", o.ReadTimeout, "ReadTimeout (in seconds) when read to single service.")
	fs.DurationVar(&o.WriteTimeout, "single.writeTimeout", o.WriteTimeout, "WriteTimeout (in seconds) when write to single service.")

}
