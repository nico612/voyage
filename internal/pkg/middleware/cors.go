// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	maxAge = 12
)

// Cors 设置跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" { // 如果是简单跨域请求，则继续处理 HTTP 请求
			c.Next()
		} else { // 复杂请求的CORS跨域处理

			// 复杂请求的 CORS 请求，会在正式通信之前，增加一次 HTTP 查询请求，称为"预检"请求（preflight）。
			// "预检"请求用的请求方法是OPTIONS，表示这个请求是用来询问请求能否安全送出的。
			//  预检通过后，浏览器就正常发起请求和响应，流程和简单请求一致。

			// 必选，设置允许访问的域名， * 接受任意域名的请求， 如果不返回这个头部，浏览器会抛出跨域错误
			c.Header("Access-Control-Allow-Origin", "*")

			// 必选，逗号分隔的字符串，表明服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")

			//逗号分隔的字符串，表明服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段。如果浏览器请求包括 Access-Control-Request-Headers 字段，则此字段是必选的
			c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")

			// Access-Control-Allow-Credentials：可选，布尔值，默认是false，表示不允许发送 Cookie
			//Access-Control-Max-Age：指定本次预检请求的有效期，单位为秒。可以避免频繁的预检请求

			c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Content-Type", "application/json")
			c.AbortWithStatus(200)
		}
	}
}

// CorsFormCorsLib 使用三方 cors 库创建
func CorsFormCorsLib() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: maxAge * time.Hour,
	})
}
