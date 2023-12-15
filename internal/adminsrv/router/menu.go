package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/adminsrv/controller/v1/menu"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/redis/go-redis/v9"
)

func InitMenuRouter(r *gin.RouterGroup, store store.IStore, rdb *redis.Client) (R gin.IRouter) {
	menuRouter := r.Group("menu")
	//.Use(middleware.OperationRecord())
	menuRouterWithoutRecord := r.Group("menu")

	menuCtr := menu.NewAuthorityMenuController(store, rdb)
	{
		menuRouter.POST("addBaseMenu", menuCtr.AddBaseMenu)           // 新增菜单
		menuRouter.POST("addMenuAuthority", menuCtr.AddMenuAuthority) //	增加menu和角色关联关系
		menuRouter.POST("deleteBaseMenu", menuCtr.DeleteBaseMenu)     // 删除菜单
		menuRouter.POST("updateBaseMenu", menuCtr.UpdateBaseMenu)     // 更新菜单
	}

	{
		menuRouterWithoutRecord.POST("getMenu", menuCtr.GetMenu)                   // 获取菜单树
		menuRouterWithoutRecord.POST("getMenuList", menuCtr.GetMenuList)           // 分页获取基础menu列表
		menuRouterWithoutRecord.POST("getBaseMenuTree", menuCtr.GetBaseMenuTree)   // 获取用户动态路由
		menuRouterWithoutRecord.POST("getMenuAuthority", menuCtr.GetMenuAuthority) // 获取指定角色menu
		menuRouterWithoutRecord.POST("getBaseMenuById", menuCtr.GetBaseMenuById)   // 根据id获取菜单
	}

	return menuRouter
}
