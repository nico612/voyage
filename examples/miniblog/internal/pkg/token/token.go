// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package token

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Config 包括 token 包的配置选项.
type Config struct {
	// 签发 token 需要用到 secret key，
	key string
	// 类似用户名这种唯一的身份标识，一般会用 identityKey 来指代，identityKey 可以保存在 Token 的 Claims 中
	// 所以，我们签发 token 时，要将用户名保存在 token 中，解析 token 时，要从 token 中解析出用户名并返回。
	identityKey string
}

// ErrMissingHeader 表示 `Authorization` 请求头为空.
var ErrMissingHeader = errors.New("the length of the `Authorization` header is zero")

var (
	// 包级别 config.
	config = Config{"Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", "identityKey"}
	// 防止同一个服务进程多次初始化，需要使用 sync.Once 来确保 token 包只被初始化一次.
	once sync.Once
)

// Init 设置包级别的配置 config, config 会用于本包后面的 token 签发和解析.
func Init(key string, identityKey string) {
	once.Do(func() {
		if key != "" {
			config.key = key
		}
		if identityKey != "" {
			config.identityKey = identityKey
		}
	})
}

// Parse 使用指定的密钥 key 解析 token，解析成功返回 token 上下文，否则报错.
func Parse(tokenString string, key string) (string, error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 通过类型断言，确保 token 加密算法是预期的加密算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}

	var identityKey string
	// 如果解析成功，从 token 中取出 token 的主题
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identityKey = claims[config.identityKey].(string)
	}
	return identityKey, nil
}

// ParseRequest 从请求头中获取令牌，并将其传递给 Parse 函数以解析令牌.
func ParseRequest(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		return "", ErrMissingHeader
	}

	var t string
	// 从请求头中取出 token
	fmt.Sscanf(header, "Bearer %s", &t)

	return Parse(t, config.key)
}

// Sign 使用 jwtSecret 签发 token，token 的 claims 中会存放传入的 subject.
func Sign(identityKey string) (tokenString string, err error) {
	// token 的内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,
		"nbf":              time.Now().Unix(),                         // JWT Token 生效时间 ；
		"iat":              time.Now().Unix(),                         // JWT 签发时间
		"exp":              time.Now().Add(100000 * time.Hour).Unix(), // 过期时间
	})

	// 签发token
	tokenString, err = token.SignedString([]byte(config.key))
	return
}
