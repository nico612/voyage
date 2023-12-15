package options

import (
	"encoding/json"
	genericoptions "github.com/nico612/voyage/internal/pkg/options"
	"github.com/nico612/voyage/internal/pkg/server"
	cliflag "github.com/nico612/voyage/pkg/cli/flag"
	"github.com/nico612/voyage/pkg/log"
	"github.com/nico612/voyage/pkg/utils/idutil"
)

// Options adminsrv 服务配置选项
type Options struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions       `json:"server"   mapstructure:"server"`   // 服务配置
	GRPCOptions             *genericoptions.GRPCOptions            `json:"grpc"     mapstructure:"grpc"`     // grpc 配置
	InsecureServing         *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"` // http 配置
	SecureServing           *genericoptions.SecureServingOptions   `json:"secure"   mapstructure:"secure"`   // https 配置
	MySQLOptions            *genericoptions.MySQLOptions           `json:"mysql"    mapstructure:"mysql"`
	RedisOptions            *genericoptions.RedisOptions           `json:"redis"    mapstructure:"redis"`
	JwtOptions              *genericoptions.JwtOptions             `json:"jwt"      mapstructure:"jwt"`
	Log                     *log.Options                           `json:"log"      mapstructure:"log"`
	FeatureOptions          *genericoptions.FeatureOptions         `json:"feature"  mapstructure:"feature"` // 特性配置
	Captcha                 *genericoptions.Captcha                `json:"captcha" mapstructure:"captcha"`  // 图形验证码
}

func NewOptions() *Options {
	o := Options{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		GRPCOptions:             genericoptions.NewGRPCOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		MySQLOptions:            genericoptions.NewMySQLOptions(),
		RedisOptions:            genericoptions.NewRedisOptions(),
		JwtOptions:              genericoptions.NewJwtOptions(),
		Log:                     log.NewOptions(),
		FeatureOptions:          genericoptions.NewFeatureOptions(),
	}

	return &o
}

// ApplyTo 应用到通用服务配置
func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

// Flags 用于返回特定 APIServer 的各个部分的标志（flags）
// 返回这些命名标志集，其中每个部分都有自己的一组标志。
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.JwtOptions.AddFlags(fss.FlagSet("jwt"))
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("single"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.Log.AddFlags(fss.FlagSet("logs"))

	return fss
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.GRPCOptions.Validate()...)
	errs = append(errs, o.InsecureServing.Validate()...)
	errs = append(errs, o.SecureServing.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.JwtOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.FeatureOptions.Validate()...)

	return errs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {
	if o.JwtOptions.Key == "" {
		o.JwtOptions.Key = idutil.NewSecretKey()
	}

	return o.SecureServing.Complete()
}
