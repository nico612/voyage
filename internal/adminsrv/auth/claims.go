package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// jwt custom claims
const (
	// APIServerIssuer 定义 jwt 签发者
	APIServerIssuer = "voyage.server.com"

	// APIServerAudience jwt 接收者
	APIServerAudience = "voyage-admin-server"

	ClaimsKey = "claims"
)

type Options struct {
	Key        string        `json:"key"`     // 密钥
	Timeout    time.Duration `json:"timeout"` // 过期时间
	MaxRefresh time.Duration `json:"max-refresh"`
}

type CustomClaims struct {
	jwt.RegisteredClaims
	MaxRefresh  time.Duration // 最大刷新时间，在该时间内(生效时间+MaxRefresh)可以刷新Token
	UserID      uint          `json:"userID"`      // userID
	Username    string        `json:"username"`    // 用户名
	AuthorityId uint          `json:"authorityId"` // 授权ID
}

func NewClaims(opts *Options) *CustomClaims {
	return &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    APIServerIssuer,                                      // 签发者
			Audience:  jwt.ClaimStrings{APIServerAudience},                  // 接收者
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1 * time.Second)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(opts.Timeout)),     // 过期时间
		},
		MaxRefresh: opts.MaxRefresh, // 刷新时间
	}
}

func GetUserID(c *gin.Context) uint {
	claims, err := GetClaims(c)
	if err != nil {
		return 0
	}
	return claims.UserID
}

func GetUserAuthorityID(c *gin.Context) uint {
	claims, err := GetClaims(c)
	if err != nil {
		return 0
	}
	return claims.AuthorityId
}

func GetClaims(c *gin.Context) (*CustomClaims, error) {
	value, exist := c.Get(ClaimsKey)
	if !exist {
		return nil, errors.New("not found claims")
	}

	claims, ok := value.(*CustomClaims)
	if !ok {
		return nil, errors.New("claims type error")
	}

	return claims, nil

}

//func (c *CustomClaims) ToMapClaims() map[string]interface{} {
//	var mapClaims map[string]interface{}
//	// 自定义扁平化解析器
//	decoder, _ := mapstructure.NewDecoder(
//		&mapstructure.DecoderConfig{
//			Squash: true,
//			Result: &mapClaims,
//		})
//	_ = decoder.Decode(c)
//
//	return mapClaims
//}
