package options

import "github.com/spf13/pflag"

// ClientCertAuthenticationOptions 定义客户端证书认证的不同选项
type ClientCertAuthenticationOptions struct {
	// ClientCA 用于验证传入客户端证书的证书捆绑包。
	ClientCA string `json:"client-ca-file" mapstructure:"client-ca-file"`
}

func NewClientCertAuthenticationOptions() *ClientCertAuthenticationOptions {
	return &ClientCertAuthenticationOptions{
		ClientCA: "",
	}
}

// Validate 验证用户在程序启动时通过命令行输入的参数。
func (o *ClientCertAuthenticationOptions) Validate() []error {
	return []error{}
}

// AddFlags 将标志添加到指定的 pflag.FlagSet.
func (o *ClientCertAuthenticationOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ClientCA, "client-ca-file", o.ClientCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the client-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")
}
