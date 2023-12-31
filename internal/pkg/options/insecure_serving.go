package options

import (
	"fmt"
	"github.com/spf13/pflag"
)

// InsecureServingOptions http 服务配置
type InsecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
}

func NewInsecureServingOptions() *InsecureServingOptions {
	return &InsecureServingOptions{
		BindAddress: "127.0.0.1",
		BindPort:    8080,
	}
}

func (i *InsecureServingOptions) Address() string {
	return fmt.Sprintf("%s:%d", i.BindAddress, i.BindPort)
}

// ApplyTo 将配置应用到 api server 配置中
//func (s *InsecureServingOptions) ApplyTo(c *server.Config) error {
//	c.InsecureServing = &server.InsecureServingInfo{
//		Address: net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort)),
//	}
//
//	return nil
//}

// Validate 验证用户在程序启动时通过命令行输入的参数。
func (s *InsecureServingOptions) Validate() []error {
	var errors []error

	if s.BindPort < 0 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--insecure.bind-port %v must be between 0 and 65535, inclusive. 0 for turning off insecure (HTTP) port",
				s.BindPort,
			),
		)
	}

	return errors
}

// AddFlags 向 APIServer 添加 标志位到 指定的 FlagSet 中
func (s *InsecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.BindAddress, "insecure.bind-address", s.BindAddress, ""+
		"The IP address on which to serve the --insecure.bind-port "+
		"(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	fs.IntVar(&s.BindPort, "insecure.bind-port", s.BindPort, ""+
		"The port on which to serve unsecured, unauthenticated access. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine and that port 443 on the iam public address is proxied to this "+
		"port. This is performed by nginx in the default setup. Set to zero to disable.")
}
