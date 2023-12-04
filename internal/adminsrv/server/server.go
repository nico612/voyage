package server

import (
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/internal/adminsrv/store/cache"
	"github.com/nico612/voyage/internal/adminsrv/store/mysql"
	genericServer "github.com/nico612/voyage/internal/pkg/server"
	"github.com/nico612/voyage/pkg/log"
	"github.com/nico612/voyage/pkg/shutdown"
	"github.com/nico612/voyage/pkg/shutdown/managers/posixsignal"
	"github.com/songzhibin97/gkit/cache/local_cache"
)

// apiServer
type apiServer struct {
	gs *shutdown.GracefulShutdown // 优雅关停
	//grpc             *grpcAPIServer
	genericApoServer *genericServer.GenericAPIServer // HTTP、HTTPS服务
	cfg              *config.Config                  // api Server 配置
}

// CreateAPIServer 创建服务
func CreateAPIServer(cfg *config.Config) (*apiServer, error) {

	// 优雅关停
	gs := shutdown.New()
	// 添加关停信号管理器
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	// http https 服务配置
	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	// 使用完整配置并新建服务
	genericServer, err := genericConfig.Complete().New()

	// TODO GRPC 服务

	// 服务扩展配置
	//extraConfig, err := buildExtraConfig(cfg)
	//if err != nil {
	//	return nil, err
	//}

	//extraServer, err := extraConfig.com

	// 初始化 mysql
	dbStore, _ := mysql.GetMySQLStoreOr(cfg.MySQLOptions)
	store.SetStore(dbStore)

	// 初始化缓存
	_ = cache.GetLocalCacheIns(local_cache.SetDefaultExpire(cfg.JwtOptions.Timeout))

	server := &apiServer{
		gs:               gs,
		genericApoServer: genericServer,
		cfg:              cfg,
	}

	return server, nil

}

// PrepareRun 主要负责运行服务前的 路由、redis 等初始化 工作
func (s *apiServer) PrepareRun() *preparedAPIServer {

	NewRouter(s.genericApoServer.Engine, s.cfg).Initializer()

	// 初始化 redis

	// 添加优雅关停回调
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		// 关闭数据库连接

		// 关闭 grpc 服务

		// 关闭 api 服务
		s.genericApoServer.Close()

		return nil
	}))

	return &preparedAPIServer{s}
}

type preparedAPIServer struct {
	*apiServer
}

// Run 准备完成后运行服务
func (s *preparedAPIServer) Run() error {
	// 运行 grpc 服务

	// 开始监听退出信号
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return s.genericApoServer.Run()

}

// 构建 http、https 服务配置
func buildGenericConfig(cfg *config.Config) (genericConfig *genericServer.Config, lastErr error) {

	genericConfig = genericServer.NewConfig() // 默认配置

	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}
