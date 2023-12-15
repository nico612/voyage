package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/redis/go-redis/v9"

	"github.com/nico612/voyage/internal/adminsrv/controller/v1/sysuser"
	"github.com/nico612/voyage/internal/adminsrv/store"
)

// InitSysUserRouter 用户路由
func InitSysUserRouter(r *gin.RouterGroup, cfg *config.Config, storeIns store.IStore, rdb *redis.Client) {

	userRouter := r.Group("user")              // 需要记录操作的路由
	userRouterWithoutRecord := r.Group("user") // 不需要操作记录的路由
	//Use(middleware.OperationRecord())

	sysUserController := sysuser.NewSysUserController(storeIns, cfg, rdb)
	{

		// 在缓冲时间段内后端 会自动 刷新token， 可跟前端协商是自动刷新还是手动刷新，如果手动刷新，后端需要设计为过期后的缓冲区
		//userRouter.POST("refreshToken", sysUserController.RefreshToken)
		userRouter.POST("admin_register", sysUserController.Register)               // 管理员注册账号
		userRouter.POST("changePassword", sysUserController.ChangePassword)         // 用户修改密码
		userRouter.POST("setUserAuthority", sysUserController.SetUserAuthority)     // 设置用户权限
		userRouter.DELETE("deleteUser", sysUserController.DeleteUser)               // 删除用户
		userRouter.PUT("setUserInfo", sysUserController.SetUserInfo)                // 设置用户信息
		userRouter.PUT("setSelfInfo", sysUserController.SetSelfInfo)                // 设置自身信息
		userRouter.POST("setUserAuthorities", sysUserController.SetUserAuthorities) // 设置用户权限组
		userRouter.POST("resetPassword", sysUserController.ResetPassword)           // 设置用户权限组
	}

	{
		userRouterWithoutRecord.POST("getUserList", sysUserController.GetUserList) // 分页获取用户列表
		userRouterWithoutRecord.GET("getUserInfo", sysUserController.GetUserInfo)  // 获取自身信息
	}
}
