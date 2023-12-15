package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/middleware"
	"github.com/nico612/voyage/internal/adminsrv/router"
	"github.com/nico612/voyage/internal/adminsrv/store/mysql"
	"github.com/nico612/voyage/internal/adminsrv/store/redis"
	generticMW "github.com/nico612/voyage/internal/pkg/middleware"
)

type Router struct {
	Prefix string // 路由前缀
	*gin.Engine
	cfg *config.Config
}

func NewRouter(engine *gin.Engine, cfg *config.Config) *Router {
	return &Router{
		Prefix: cfg.GenericServerRunOptions.RouterPrefix,
		Engine: engine,
		cfg:    cfg,
	}
}

func (r *Router) Initializer() {

	// store 实例在服务创建时已初始化，这里直接获取就行
	storeIns, _ := mysql.GetMySQLStoreOr(nil)

	// redis
	rdb, _ := redis.CreateRedisOr(nil)

	r.Use(generticMW.Cors())

	// 不需要认证的路由
	publicRouter := r.Group(r.Prefix)
	{
		router.InstallBaseRouter(publicRouter, r.cfg, storeIns, rdb)
	}

	// jwt 中间件
	//authz := auth.NewJWTAuth().MiddlewareFunc()
	// 需要认证的路由
	privateRouter := r.Group(r.Prefix)

	privateRouter.Use(middleware.JWTAuth())
	{
		router.InitSysUserRouter(privateRouter, r.cfg, storeIns, rdb) // 注册用户路由
		router.InitAuthorityRouter(privateRouter, storeIns, rdb)
		router.InitMenuRouter(privateRouter, storeIns, rdb) // 注册menu路由
	}

}
