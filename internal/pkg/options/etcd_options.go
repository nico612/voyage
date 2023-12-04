// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/spf13/pflag"
)

// EtcdOptions 定义 etcd 集群配置选项
type EtcdOptions struct {
	Endpoints            []string `json:"endpoints"               mapstructure:"endpoints"`               // etcd 实例的端点地址列表。
	Timeout              int      `json:"timeout"                 mapstructure:"timeout"`                 // 操作的超时时间。
	RequestTimeout       int      `json:"request-timeout"         mapstructure:"request-timeout"`         // 请求超时时间
	LeaseExpire          int      `json:"lease-expire"            mapstructure:"lease-expire"`            // 租约到期时间
	Username             string   `json:"username"                mapstructure:"username"`                // 连接用户名
	Password             string   `json:"password"                mapstructure:"password"`                // 连接密码
	UseTLS               bool     `json:"use-tls"                 mapstructure:"use-tls"`                 // 是否使用 TLS 连接
	CaCert               string   `json:"ca-cert"                 mapstructure:"ca-cert"`                 // CA 证书
	Cert                 string   `json:"cert"                    mapstructure:"cert"`                    // 客户端证书
	Key                  string   `json:"key"                     mapstructure:"key"`                     // 客户端证书 私钥
	HealthBeatPathPrefix string   `json:"health_beat_path_prefix" mapstructure:"health_beat_path_prefix"` // 健康检查心跳路径前缀
	HealthBeatIFaceName  string   `json:"health_beat_iface_name"  mapstructure:"health_beat_iface_name"`  // 健康检查心跳接口名称
	Namespace            string   `json:"namespace"               mapstructure:"namespace"`               // etch 的命令空间
}

// NewEtcdOptions create a `zero` value instance.
func NewEtcdOptions() *EtcdOptions {
	return &EtcdOptions{
		Timeout:        5,
		RequestTimeout: 2,
		LeaseExpire:    5,
	}
}

// Validate verifies flags passed to RedisOptions.
func (o *EtcdOptions) Validate() []error {
	errs := []error{}

	if len(o.Endpoints) == 0 {
		errs = append(errs, fmt.Errorf("etcd endpoints can not be empty"))
	}

	if o.RequestTimeout <= 0 {
		errs = append(errs, fmt.Errorf("--etcd.request-timeout cannot be negative"))
	}

	return errs
}

// AddFlags adds flags related to redis storage for a specific APIServer to the specified FlagSet.
func (o *EtcdOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.Endpoints, "etcd.endpoints", o.Endpoints, "Endpoints of etcd cluster.")
	fs.StringVar(&o.Username, "etcd.username", o.Username, "Username of etcd cluster.")
	fs.StringVar(&o.Password, "etcd.password", o.Password, "Password of etcd cluster.")
	fs.IntVar(&o.Timeout, "etcd.timeout", o.Timeout, "Etcd dial timeout in seconds.")
	fs.IntVar(&o.RequestTimeout, "etcd.request-timeout", o.RequestTimeout, "Etcd request timeout in seconds.")
	fs.IntVar(&o.LeaseExpire, "etcd.lease-expire", o.LeaseExpire, "Etcd expire timeout in seconds.")
	fs.BoolVar(&o.UseTLS, "etcd.use-tls", o.UseTLS, "Use tls transport to connect etcd cluster.")
	fs.StringVar(&o.CaCert, "etcd.ca-cert", o.CaCert, "Path to cacert for connecting to etcd cluster.")
	fs.StringVar(&o.Cert, "etcd.cert", o.Cert, "Path to cert file for connecting to etcd cluster.")
	fs.StringVar(&o.Key, "etcd.key", o.Key, "Path to key file for connecting to etcd cluster.")
	fs.StringVar(
		&o.HealthBeatPathPrefix,
		"etcd.health-beat-path-pre",
		o.HealthBeatPathPrefix,
		"health beat path prefix.",
	)
	fs.StringVar(
		&o.HealthBeatIFaceName,
		"etcd.health-beat-iface-name",
		o.HealthBeatIFaceName,
		"health beat registry iface name, such as eth0.",
	)
	fs.StringVar(&o.Namespace, "etcd.namespace", o.Namespace, "Etcd storage namespace.")
}

// GetEtcdTLSConfig returns etcd tls config.
func (o *EtcdOptions) GetEtcdTLSConfig() (*tls.Config, error) {
	var (
		cert       tls.Certificate
		certLoaded bool
		capool     *x509.CertPool
	)
	if o.Cert != "" && o.Key != "" {
		var err error
		cert, err = tls.LoadX509KeyPair(o.Cert, o.Key)
		if err != nil {
			return nil, err
		}
		certLoaded = true
		o.UseTLS = true
	}
	if o.CaCert != "" {
		data, err := ioutil.ReadFile(o.CaCert)
		if err != nil {
			return nil, err
		}
		capool = x509.NewCertPool()
		for {
			var block *pem.Block
			block, _ = pem.Decode(data)
			if block == nil {
				break
			}
			cacert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			capool.AddCert(cacert)
		}
		o.UseTLS = true
	}

	if o.UseTLS {
		cfg := &tls.Config{
			RootCAs:            capool,
			InsecureSkipVerify: false,
		}
		if certLoaded {
			cfg.Certificates = []tls.Certificate{cert}
		}

		return cfg, nil
	}

	return &tls.Config{}, nil
}
