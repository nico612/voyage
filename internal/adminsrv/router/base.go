package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/controller/v1/base"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/redis/go-redis/v9"
)

// InstallBaseRouter  初始化 base 控制层
func InstallBaseRouter(r *gin.RouterGroup, cfg *config.Config, storeIns store.IStore, rdb *redis.Client) {

	baseRouter := r.Group("base")
	{
		basectr := base.NewBaseController(storeIns, rdb, cfg)

		baseRouter.POST("login", basectr.Login)
		baseRouter.POST("captcha", basectr.Captcha)
	}

}
