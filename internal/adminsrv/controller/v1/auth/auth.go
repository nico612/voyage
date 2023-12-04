package auth

import (
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	v1 "github.com/nico612/voyage/api/adminsrv/v1"
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/pkg/middleware"
	jwt2 "gopkg.in/dgrijalva/jwt-go.v3"
	"time"
)

const (
	// APIServerAudience 定义 jwt 签发者
	APIServerAudience = "voyage.server.com"

	// APIServerIssuer jwt 接收者
	APIServerIssuer = "voyage-admin-server"
)

type JWTMiddleWare struct {
	*jwt.GinJWTMiddleware
}

type Claims struct {
	UserID   uint
	UserName string
}

func NewJWTAuth() *JWTMiddleWare {

	opt := config.GetAppConfig().JwtOptions

	return &JWTMiddleWare{&jwt.GinJWTMiddleware{
		Realm:            opt.Realm,
		SigningAlgorithm: "HS256",
		Key:              []byte(opt.Key),
		Timeout:          opt.Timeout,
		MaxRefresh:       opt.MaxRefresh,
		Authenticator:    authenticator,
		LoginResponse:    loginResponse,
		RefreshResponse:  refreshResponse,
		PayloadFunc:      payloadFunc,
		IdentityHandler: func(claims jwt2.MapClaims) interface{} {
			// 这里jwt.IdentityKey 就为设置的 middleware.UserIDKey
			// 取出的数据为 user.UserID, 在下面 payloadFunc 回调方法中设置
			return claims[middleware.UserIDKey]
		},
		Authorizator: authorizator, // 授权回调
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		SendCookie:    true,
		TimeFunc:      time.Now,
		//HTTPStatusMessageFunc: func(e errors, c *gin.Context) string {
		//
		//},
	}}
}

// 自定义认证，登录请求进来，会回调该方法，进行登录信息验证，比如验证用户名和和密码
// 然后返回用户和错误信息
func authenticator(c *gin.Context) (interface{}, error) {
	var login v1.LoginReq
	var err error

	login, err = parseWithBody(c)
	if err != nil {
		return "", jwt.ErrFailedAuthentication
	}

	// 验证图形验证码

	// 验证用户名和密码

	// 返回用户信息

	return nil, nil
}

func parseWithBody(c *gin.Context) (v1.LoginReq, error) {

}

// 在登录成功后，后面的操作中进行授权验证
// data 为 解析token后 从 claims[jwt.IdentityKey] 取出，在 payloadFunc 方法中设置为 claims[jwt.IdentityKey] = u.UserID 用户ID
// 因此这里的 data 为 UserID，可以在该方法中做一些其他授权类的操作
func authorizator(data interface{}, c *gin.Context) bool {
	return true
}

// 设置 token 中第二段 payload 包含的信息
// data 是上面 authenticator 认证时返回的数据，这里是用户信息的数据
// 标准字段有：
// - `iss`：JWT Token 的签发者，其值大小写敏感字符串或者URL；
// - `sub`：主题；可以用来鉴别一个用户
// - `exp`：JWT Token 过期时间；
// - `aud`：接收 JWT Token 的一方，其值应为大小写敏感字符串或URL，一般可以为特定的App，服务或模块。服务端的安全策略在签发时和验证时，aud必须时一致的；
// - `iat`：JWT Token 签发时间；
// - `nbf`：JWT Token 生效时间 ；
// - `jti`：JWT Token ID。令牌唯一标识符，通常用于一次性消费的Token
func payloadFunc(data interface{}) jwt.MapClaims {
	claims := jwt.MapClaims{
		"iss": APIServerIssuer,
		"aud": APIServerAudience,
	}

	// payload 添加其他信息
	if u, ok := data.(models.SysUser); ok {
		claims[middleware.UserIDKey] = u.UserID
		claims["sub"] = u.UserID // 接收者信息
	}
	return claims
}

// 登录后回调
func loginResponse(*gin.Context, int, string, time.Time) {

}

// 刷新 token 回调
func refreshResponse(*gin.Context, int, string, time.Time) {

}
