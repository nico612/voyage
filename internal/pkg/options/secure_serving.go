package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"net"
	"path"
)

// SecureServingOptions HTTPS 服务配置
type SecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	// BindPort 在设置了 Listener 时，BindPort 字段将被忽略，即使设置为 0 也会提供 HTTPS 服务。
	// 这意味着，当已设置监听器时，BindPort 字段将不再是确定 HTTPS 服务的关键因素，服务会在指定的监听器上提供安全的 HTTPS 连接，而不是依赖于 BindPort 的值
	BindPort int `json:"bind-port"    mapstructure:"bind-port"`
	// Required 如果设置为 true，则表示 BindPort 不能为零。
	Required bool
	// ServerCert 用于存储 TLS 证书信息
	ServerCert GeneratableKeyCert `json:"tls"          mapstructure:"tls"`
	// AdvertiseAddress net.IP
}

// CertKey HTTPS 证书文件路径相关配置
type CertKey struct {
	// CertFile 包含 PEM 编码的证书文件的路径。这个文件可能包含完整的证书链。
	CertFile string `json:"cert-file"        mapstructure:"cert-file"`
	// KeyFile 包含 PEM 编码的私钥文件的路径，对应于 CertFile 指定的证书。
	KeyFile string `json:"private-key-file" mapstructure:"private-key-file"`
}

// GeneratableKeyCert 包含的与证书相关的配置项。
type GeneratableKeyCert struct {

	// CertKey 允许设置明确的证书/密钥文件。
	CertKey CertKey `json:"cert-key" mapstructure:"cert-key"`

	// CertDirectory 如果未明确设置 CertFile 和 KeyFile，这些字段允许你指定用于写入生成证书的目录以及确定在该目录中生成的证书文件的文件名。
	// 如果 CertDirectory 和 PairName 都未设置，系统将生成一个存储在内存中的证书。
	CertDirectory string `json:"cert-dir"  mapstructure:"cert-dir"`

	// PairName 证书和密钥文件名的名称
	// 将与 CertDirectory 结合使用，用于确定生成的证书和密钥文件的文件名。
	// 例如，如果 PairName 设置为 "example", 那么生成的证书文件名将是 CertDirectory/example.pem 和 CertDirectory/example.key。
	PairName string `json:"pair-name" mapstructure:"pair-name"`
}

// NewSecureServingOptions 默认配置项
func NewSecureServingOptions() *SecureServingOptions {
	return &SecureServingOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8443,
		Required:    true,
		ServerCert: GeneratableKeyCert{
			PairName:      "adminsrv",
			CertDirectory: "/var/run/adminsrv",
		},
	}
}

//// ApplyTo 将配置信息引用到 api server 配置中
//func (s *SecureServingOptions) ApplyTo(c *server.Config) error {
//	// SecureServing is required to serve https
//	c.SecureServing = &server.SecureServingInfo{
//		BindAddress: s.BindAddress,
//		BindPort:    s.BindPort,
//		CertKey: server.CertKey{
//			CertFile: s.ServerCert.CertKey.CertFile,
//			KeyFile:  s.ServerCert.CertKey.KeyFile,
//		},
//	}
//
//	return nil
//}

// Validate 验证用户在程序启动时通过命令行输入的参数。
func (s *SecureServingOptions) Validate() []error {
	if s == nil {
		return nil
	}

	errors := []error{}

	if s.Required && s.BindPort < 1 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--secure.bind-port %v must be between 1 and 65535, inclusive. It cannot be turned off with 0",
				s.BindPort,
			),
		)
	} else if s.BindPort < 0 || s.BindPort > 65535 {
		errors = append(errors, fmt.Errorf("--secure.bind-port %v must be between 0 and 65535, inclusive. 0 for turning off secure port", s.BindPort))
	}

	return errors
}

// AddFlags 向 APIServer 添加 标志位到 指定的 FlagSet 中
func (s *SecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.BindAddress, "secure.bind-address", s.BindAddress, ""+
		"The IP address on which to listen for the --secure.bind-port port. The "+
		"associated interface(s) must be reachable by the rest of the engine, and by CLI/web "+
		"clients. If blank, all interfaces will be used (0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	desc := "The port on which to serve HTTPS with authentication and authorization."
	if s.Required {
		desc += " It cannot be switched off with 0."
	} else {
		desc += " If 0, don't serve HTTPS at all."
	}
	fs.IntVar(&s.BindPort, "secure.bind-port", s.BindPort, desc)

	fs.StringVar(&s.ServerCert.CertDirectory, "secure.tls.cert-dir", s.ServerCert.CertDirectory, ""+
		"The directory where the TLS certs are located. "+
		"If --secure.tls.cert-key.cert-file and --secure.tls.cert-key.private-key-file are provided, "+
		"this flag will be ignored.")

	fs.StringVar(&s.ServerCert.PairName, "secure.tls.pair-name", s.ServerCert.PairName, ""+
		"The name which will be used with --secure.tls.cert-dir to make a cert and key filenames. "+
		"It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key")

	fs.StringVar(&s.ServerCert.CertKey.CertFile, "secure.tls.cert-key.cert-file", s.ServerCert.CertKey.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")

	fs.StringVar(&s.ServerCert.CertKey.KeyFile, "secure.tls.cert-key.private-key-file",
		s.ServerCert.CertKey.KeyFile, ""+
			"File containing the default x509 private key matching --secure.tls.cert-key.cert-file.")
}

// Complete 确保 SecureServingOptions 结构中涉及安全服务的配置项被完整设置。
func (s *SecureServingOptions) Complete() error {
	if s == nil || s.BindPort == 0 {
		return nil
	}

	keyCert := &s.ServerCert.CertKey
	if len(keyCert.CertFile) != 0 || len(keyCert.KeyFile) != 0 {
		return nil
	}

	if len(s.ServerCert.CertDirectory) > 0 {
		if len(s.ServerCert.PairName) == 0 {
			return fmt.Errorf("--secure.tls.pair-name is required if --secure.tls.cert-dir is set")
		}
		keyCert.CertFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".crt")
		keyCert.KeyFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".key")
	}

	return nil
}

// CreateListener 通过给定的 address 创建 listener 并返回 listener 和 端口.
func CreateListener(addr string) (net.Listener, int, error) {
	network := "tcp"

	ln, err := net.Listen(network, addr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to listen on %v: %w", addr, err)
	}

	// get port
	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		_ = ln.Close()

		return nil, 0, fmt.Errorf("invalid listen address: %q", ln.Addr().String())
	}

	return ln, tcpAddr.Port, nil
}
