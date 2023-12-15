package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/nico612/voyage/pkg/log"
	"golang.org/x/sync/singleflight"
	"strings"
)

type JWTAuth struct {
	// 用于签名Secret key. Required
	Key []byte

	// 签名算法 - 可用 HS256, HS384, HS5512, RS256, RS384 or RS512
	//Optional 默认 HS256
	SigningAlgorithm string

	// 在并发调用方法时，使用归并回源避免并发问题
	group *singleflight.Group
}

func NewJwtAuth(key []byte) *JWTAuth {
	return &JWTAuth{
		Key:              key,
		SigningAlgorithm: "HS256",
		group:            &singleflight.Group{},
	}
}

// GeneratorToken 生成token
func (g *JWTAuth) GeneratorToken(claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod(g.SigningAlgorithm), claims)
	return token.SignedString(g.Key)
}

// ParseTokenString 解析token
func (g *JWTAuth) ParseTokenString(tokenString string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return g.Key, nil
	})

	if err != nil {
		return nil, err
	}
	//if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
	//	return claims, nil
	//}
	claims := token.Claims.(*CustomClaims)
	return claims, nil
}

// RefreshToken 刷新token
func (g *JWTAuth) RefreshToken(old string, claims *CustomClaims) (string, error) {

	// 使用归并回源避免并发问题
	v, err, _ := g.group.Do("JWT:"+old, func() (interface{}, error) {
		return g.GeneratorToken(claims)
	})
	return v.(string), err
}

// CheckIfTokenExpire 检查是否过期，如没有过期返回解析出来的信息
func (g *JWTAuth) CheckIfTokenExpire(tokenStr string) (*CustomClaims, error) {

	// 判断是否在刷新时间范围内
	claims, err := g.ParseTokenString(tokenStr)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func ParseJwtTokenFormHeader(c *gin.Context) (string, error) {
	// jwt token 在 header 中的 格式为 Authorization : Bearer xxxxx.xxxxx.xxxxxx
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Bearer" {
		log.Errorf("get basic string from Authorization header failed")
		return "", errors.New("failed authentication")
	}
	return auth[1], nil
}

func ParseJwtTokenFormApiKey(c *gin.Context) (string, error) {
	token := c.Request.Header.Get("x-token")
	if token == "" {
		return "", errors.New("token is empty")
	}
	return token, nil
}
