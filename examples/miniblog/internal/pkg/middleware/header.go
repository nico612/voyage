package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// NoCache 是一个 Gin 中间件，用来禁止客户端缓存 HTTP 请求的返回结果.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Cors 是一个 Gin 中间件，用来设置 options 请求的返回头，然后退出中间件链，并结束请求(浏览器跨域设置).
// 复杂请求的 CORS 跨域处理
// 复杂请求的 CORS 请求，会在正式通信之前，增加一次 HTTP 查询请求，称为"预检"请求（preflight）。"预检"请求用的请求方法是OPTIONS，表示这个请求是用来询问请求能否安全送出的。
// 预检通过后，浏览器就正常发起请求和响应，流程和简单请求一致。
// 当后端收到预检请求后，可以设置跨域相关 Header 以完成跨域请求。支持的 Header 具体如下表所示：
func Cors(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else { // 预检：发送OPTIONS请求。设置跨域Header，并返回

		//必选，表示接受任意域名的请求。如果不返回这个头部，浏览器会抛出跨域错误。（注：服务器端不会报错）
		c.Header("Access-Control-Allow-Origin", "*")
		// 必选，逗号分隔的字符串，表明服务器支持的所有跨域请求的方法
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		// 逗号分隔的字符串，表明服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段。如果浏览器请求包括 Access-Control-Request-Headers 字段，则此字段是必选的
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		//Access-Control-Allow-Credentials	可选，布尔值，默认是false，表示不允许发送 Cookie
		//c.Header("Access-Control-Allow-Credentials", "false")

		//指定本次预检请求的有效期，单位为秒。可以避免频繁的预检请求
		//c.Header("Access-Control-Max-Age	", "100")

		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(200)
	}
}

// Secure 是一个 Gin 中间件，用来添加一些安全和资源访问相关的 HTTP 头.
func Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
}
