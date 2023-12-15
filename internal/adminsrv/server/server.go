package server

import (
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/internal/adminsrv/store/cache"
	"github.com/nico612/voyage/internal/adminsrv/store/mysql"
	"github.com/nico612/voyage/internal/adminsrv/store/redis"
	genericServer "github.com/nico612/voyage/internal/pkg/server"
	"github.com/nico612/voyage/internal/pkg/utils/validator"
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

	// 初始化 mysql
	dbStore, err := mysql.GetMySQLStoreOr(cfg.MySQLOptions)
	if err != nil {
		log.Fatalf("Failed to get cache instance: %s", err.Error())
		return nil, err
	}
	store.SetClient(dbStore)

	// redis
	_, err = redis.CreateRedisOr(cfg.RedisOptions)
	if err != nil {
		log.Fatalf("Failed to get cache instance: %s", err.Error())
		return nil, err
	}

	// 初始化缓存, 并加载所有的 黑名单 jwt token 到本地缓存
	localCache := cache.GetLocalCacheIns(local_cache.SetDefaultExpire(cfg.JwtOptions.Timeout))
	if err = localCache.LoadAllJwtBlackList(dbStore); err != nil {
		log.Errorf("load all jwt black list error = %s", err.Error())
	}

	// TODO 初始化参数验证器， 不太好用需要修改
	validator.Initialize()

	server := &apiServer{
		gs:               gs,
		genericApoServer: genericServer,
		cfg:              cfg,
	}

	return server, nil

}

// PrepareRun 主要负责运行服务前的 路由、single 等初始化 工作
func (s *apiServer) PrepareRun() *preparedAPIServer {

	// 路由
	NewRouter(s.genericApoServer.Engine, s.cfg).Initializer()

	// 添加优雅关停回调
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		// 关闭数据库连接
		mysqlstore, _ := mysql.GetMySQLStoreOr(nil)
		if mysqlstore != nil {
			_ = mysqlstore.Close()
		}

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

	genericConfig = &genericServer.Config{
		InsecureServing: &genericServer.InsecureServingInfo{
			Address: cfg.InsecureServing.Address(),
		},
		Mode:            cfg.GenericServerRunOptions.Mode,
		Healthz:         cfg.GenericServerRunOptions.Healthz,
		EnableProfiling: true,
		EnableMetrics:   true,
	}
	return
}
