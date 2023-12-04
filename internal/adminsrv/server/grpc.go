package server

import (
	"fmt"
	"github.com/nico612/voyage/internal/adminsrv/config"
	genericoptions "github.com/nico612/voyage/internal/pkg/options"
	"github.com/nico612/voyage/pkg/log"
	"google.golang.org/grpc"
	"net"
)

// grpc api server
type grpcAPIServer struct {
	*grpc.Server
	address string
}

// grpc api server 配置
type grpcConfig struct {
	//Addr         string
	//MaxMsgSize   int
	//ServerCert   genericoptions.GeneratableKeyCert
	*genericoptions.GRPCOptions
	mysqlOptions *genericoptions.MySQLOptions
	// etcdOptions      *genericoptions.EtcdOptions
}

func buildGrpcConfig(cfg *config.Config) (*grpcConfig, error) {
	return &grpcConfig{
		GRPCOptions:  cfg.GRPCOptions,
		mysqlOptions: cfg.MySQLOptions,
		// etcdOptions:      cfg.EtcdOptions,
	}, nil
}

func (g *grpcConfig) address() string {
	return fmt.Sprintf("%s:%d", g.BindAddress, g.BindPort)
}

func (g *grpcConfig) complete() (grpcCompleteConfig, error) {
	return grpcCompleteConfig{g}, nil
}

// ExtraConfig 定义 grpc 完善配置
type grpcCompleteConfig struct {
	*grpcConfig
}

func (g *grpcCompleteConfig) New() (*grpcAPIServer, error) {
	//creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
	//if err != nil {
	//	log.Fatalf("Failed to generate credentials %s", err.Error())
	//}

	//opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize), grpc.Creds(creds)}
	//opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize)}

	// 创建 grpc 服务
	//grpcServer := grpc.NewServer(opts...)

	// 初始化数据库相关实例

	// 注册 grpc api 实现服务

	return nil, nil

}

// nolint: unparam
// buildExtraConfig 构建扩展配置
//func buildExtraConfig(cfg *config.Config) (*ExtraConfig, errors) {
//	return &ExtraConfig{
//		Addr:         fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
//		MaxMsgSize:   cfg.GRPCOptions.MaxMsgSize,
//		ServerCert:   cfg.SecureServing.ServerCert,
//		mysqlOptions: cfg.MySQLOptions,
//		// etcdOptions:      cfg.EtcdOptions,
//	}, nil
//}

// 完整扩展配置
//type completedExtraConfig struct {
//	*ExtraConfig
//}
//
//// 对配置进行填充或其他操作后返回一个完整的扩展配置
//func (c *ExtraConfig) complete() *completedExtraConfig {
//	if c.Addr == "" {
//		c.Addr = "127.0.0.1:8081"
//	}
//
//	return &completedExtraConfig{c}
//}

// New 创建 grpcAPIServer 实例
//func (c *completedExtraConfig) New() (*grpcAPIServer, errors) {
//	creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
//	if err != nil {
//		log.Fatalf("Failed to generate credentials %s", err.Error())
//	}
//	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize), grpc.Creds(creds)}
//	grpcServer := grpc.NewServer(opts...)
//
//	//storeIns, _ := mysql.GetMySQLFactoryOr(c.mysqlOptions)
//	//// storeIns, _ := etcd.GetEtcdFactoryOr(c.etcdOptions, nil)
//	//store.SetClient(storeIns)
//	//cacheIns, err := cachev1.GetCacheInsOr(storeIns)
//	//if err != nil {
//	//	log.Fatalf("Failed to get cache instance: %s", err.Error())
//	//}
//	//
//	//pb.RegisterCacheServer(grpcServer, cacheIns)
//	//
//	//reflection.Register(grpcServer)
//
//	return &grpcAPIServer{grpcServer, c.Addr}, nil
//}

func (s *grpcAPIServer) Run() {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}

	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalf("failed to start grpc server: %s", err.Error())
		}
	}()

	log.Infof("start grpc server at %s", s.address)
}

func (s *grpcAPIServer) Close() {
	s.GracefulStop()
	log.Infof("GRPC server on %s stopped", s.address)
}
