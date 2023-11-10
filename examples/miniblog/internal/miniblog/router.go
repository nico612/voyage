package miniblog

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/controller/v1/user"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/store"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/auth"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/core"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/errno"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/log"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/middleware"
)

// installRouters 安装 miniblog 接口路由.
func installRouters(g *gin.Engine) error {

	// 注册 404 Handle
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	// 注册 /healthz handle
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")
		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	authz, err := auth.NewAuthz(store.S.DB())
	if err != nil {
		return err
	}

	uc := user.New(store.S, authz)

	g.POST("/login", uc.Login)

	v1 := g.Group("/v1")
	{
		// 创建 users 路由分组
		userV1 := v1.Group("/users")

		{
			userV1.POST("", uc.Create)                              // 创建用户
			userV1.PUT(":name/change-password", uc.ChangePassword)  // 修改用户密码
			userV1.Use(middleware.Authn(), middleware.Authz(authz)) //认证要在授权之前，因为只有通过认证，我们才能获取到用户名，并将用户名添加到 *gin.Context 中，供授权时使用。
			userV1.GET(":name", uc.Get)
		}
	}
	return nil
}
