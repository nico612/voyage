package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/adminsrv/controller/v1/authority"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/redis/go-redis/v9"
)

func InitAuthorityRouter(r *gin.RouterGroup, store store.IStore, rdb *redis.Client) {
	authorityRouter := r.Group("authority")
	//.Use(middleware.OperationRecord())
	authorityRouterWithoutRecord := r.Group("authority")
	authorityCtr := authority.NewAuthorityController(store, rdb)

	{
		authorityRouter.POST("createAuthority", authorityCtr.CreateAuthority)   // 创建角色
		authorityRouter.POST("deleteAuthority", authorityCtr.DeleteAuthority)   // 删除角色
		authorityRouter.PUT("updateAuthority", authorityCtr.UpdateAuthority)    // 更新角色
		authorityRouter.POST("copyAuthority", authorityCtr.CopyAuthority)       // 拷贝角色
		authorityRouter.POST("setDataAuthority", authorityCtr.SetDataAuthority) // 设置角色资源权限
	}
	{
		authorityRouterWithoutRecord.POST("getAuthorityList", authorityCtr.GetAuthorityList) // 获取角色列表
	}
}
