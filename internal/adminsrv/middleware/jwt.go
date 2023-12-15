package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nico612/voyage/internal/adminsrv/auth"
	"github.com/nico612/voyage/internal/adminsrv/config"

	//"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/internal/pkg/response"
	"github.com/nico612/voyage/pkg/errors"
	"strconv"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		token, err := auth.ParseJwtTokenFormApiKey(c)
		if err != nil {
			response.Failed(c, err)
			c.Abort()
			return
		}

		jwtOpts := config.GetAppConfig().JwtOptions

		// 解析token
		jwtAuth := auth.NewJwtAuth([]byte(jwtOpts.Key))
		claims, err := jwtAuth.ParseTokenString(token)
		if err != nil {
			// token 过期
			if errors.Is(err, jwt.ErrTokenExpired) {
				response.Failed(c, errors.WithCode(code.ErrTokenExpired, "token expired"))
				c.Abort()
				return
			}

			// 无效的token
			response.Failed(c, errors.WithCode(code.ErrTokenInvalid, "token invalid"))
			c.Abort()
			return
		}

		// 在缓冲区刷新token
		if claims.ExpiresAt.Unix()-time.Now().Unix() < int64(claims.MaxRefresh) {

			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(jwtOpts.Timeout))
			newToken, _ := jwtAuth.RefreshToken(token, claims)
			newClaims, _ := jwtAuth.ParseTokenString(newToken)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))
			claims = newClaims
			// TODO 处理单点登录
		}

		// 将 UserId 设置到上下文中
		c.Set(auth.ClaimsKey, claims)
		c.Next()
	}
}
