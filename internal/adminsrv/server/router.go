package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/router"
)

type Router struct {
	Prefix string // 路由前缀
	*gin.Engine
	cfg *config.Config
}

func NewRouter(engine *gin.Engine, cfg *config.Config) *Router {
	return &Router{
		Prefix: "/adminsrv",
		Engine: engine,
		cfg:    cfg,
	}
}

func (r *Router) Initializer() {
	r.initPublicRouter()
	r.initPrivateRouter()
}

func (r *Router) initPublicRouter() {
	publicRouter := r.Group(r.Prefix)
	{
		router.InitBaseRouter(publicRouter, r.cfg)
	}
}

func (r *Router) initPrivateRouter() {
	//privateRouter := r.Group(r.Prefix)

}
