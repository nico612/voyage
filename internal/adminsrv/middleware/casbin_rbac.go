package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/adminsrv/auth"
	mycasbin "github.com/nico612/voyage/internal/adminsrv/casbin"
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/internal/pkg/response"
	"github.com/nico612/voyage/pkg/errors"
	"github.com/nico612/voyage/pkg/log"
	"strconv"
	"strings"
)

// CashbinHandler api 拦截
func CashbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := config.GetAppConfig()
		if config.GenericServerRunOptions.Mode != "debug" {

			token, err := auth.ParseJwtTokenFormApiKey(c)
			if err != nil {
				log.Errorf("casbin parse token err: %s", err.Error())
				c.Abort()
				return
			}

			jwt := auth.NewJwtAuth([]byte(config.JwtOptions.Key))
			waitUse, err := jwt.ParseTokenString(token)
			if err != nil {
				log.Errorf("cashbin parse token string err := %s", err.Error())
				c.Abort()
				return
			}
			path := c.Request.URL.Path
			// 请求实体
			sub := strconv.Itoa(int(waitUse.AuthorityId))
			// 资源
			obj := strings.TrimPrefix(path, config.GenericServerRunOptions.RouterPrefix)
			// 方法
			act := c.Request.Method
			e := mycasbin.CreateCasbinOr(nil)
			success, _ := e.Enforce(sub, obj, act)
			if !success {
				response.Failed(c, errors.WithCode(code.ErrInsufficientPermissions, "权限不足"))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
