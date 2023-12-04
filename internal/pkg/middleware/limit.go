package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ErrLimitExceeded defines Limit exceeded errors.
var ErrLimitExceeded = errors.New("Limit exceeded")

// Limit 限制请求的频率。当请求到达限制 报错 (HTTP status 429)
// maxEventsPerSec 最大每秒事件数
// maxBurstSize 最大突发大小作为参数
func Limit(maxEventsPerSec float64, maxBurstSize int) gin.HandlerFunc {

	limiter := rate.NewLimiter(rate.Limit(maxEventsPerSec), maxBurstSize)

	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()

			return
		}

		// Limit reached
		_ = c.Error(ErrLimitExceeded)
		c.AbortWithStatus(429)
	}
}
