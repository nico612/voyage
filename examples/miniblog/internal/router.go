package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/controller/v1/user"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/store"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/core"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/errno"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/log"
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

	uc := user.New(store.S)
	v1 := g.Group("/v1")
	{
		// 创建 users 路由分组
		userV1 := v1.Group("/users")
		{
			userV1.POST("/create", uc.Create)
		}
	}
	return nil
}
