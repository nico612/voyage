package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/controller/v1/base"
	"github.com/nico612/voyage/internal/adminsrv/store/mysql"
)

type BaseRouter struct {
}

func InitBaseRouter(r *gin.RouterGroup, cfg *config.Config) {

	// store 实例在服务创建时已初始化，这里直接获取就行
	storeIns, _ := mysql.GetMySQLStoreOr(nil)

	baseRouter := r.Group("base")
	{
		basectr := base.NewBaseController(storeIns, cfg)

		baseRouter.POST("login", basectr.Login)
		baseRouter.POST("captcha", basectr.Captcha)
	}

}
