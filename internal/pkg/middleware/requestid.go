package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io"
	"time"
)

const (
	// XRequestIDKey 请求头定义请求ID
	XRequestIDKey = "X-Request-ID"
)

// RequestID z中间件，用于注入 "X-Request-ID" 到上下文 和 request/response header 中
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(XRequestIDKey)

		if rid == "" {
			rid = uuid.NewV4().String()
			c.Request.Header.Set(XRequestIDKey, rid)
		}

		// Set XRequestIDKey header
		c.Writer.Header().Set(XRequestIDKey, rid)
		c.Next()
	}
}

// GetLoggerConfig 自定义 gin.LoggerConfig ， 该配置通过给定的 gin.LoFormatter 指定日志输出路径
// gin 日志 默认输出 gin.DefaultWriter = os.Stdout
// reference: https://github.com/gin-gonic/gin#custom-log-format
func GetLoggerConfig(formatter gin.LogFormatter, output io.Writer, skipPaths []string) gin.LoggerConfig {
	if formatter == nil {
		formatter = GetDefaultLogFormatterWithRequestID()
	}

	return gin.LoggerConfig{
		Formatter: formatter,
		Output:    output,
		SkipPaths: skipPaths,
	}
}

// GetDefaultLogFormatterWithRequestID returns gin.LogFormatter with 'RequestID'.
func GetDefaultLogFormatterWithRequestID() gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			param.Latency -= param.Latency % time.Second
		}

		// 返回日志格式化输出内容
		return fmt.Sprintf("%s%3d%s - [%s] \"%v %s%s%s %s\" %s",
			// param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.ClientIP,
			param.Latency,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}
}

// GetRequestIDFromContext 从给定的 context 上下文中 返回 'RequestID'
func GetRequestIDFromContext(c *gin.Context) string {
	if v, ok := c.Get(XRequestIDKey); ok {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}

	return ""
}

// GetRequestIDFromHeaders 从请求 header 中返回 'RequestID'
func GetRequestIDFromHeaders(c *gin.Context) string {
	return c.Request.Header.Get(XRequestIDKey)
}
