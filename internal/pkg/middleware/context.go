// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/internal/adminsrv/auth"
	"github.com/nico612/voyage/pkg/log"
)

// UsernameKey defines the key in gin context which represents the owner of the secret.
const UsernameKey = "username"
const UserIDKey = "userID"

// Context is a middleware that injects common prefix fields to gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 注入上下文
		c.Set(log.KeyRequestID, c.GetString(XRequestIDKey))
		c.Set(log.KeyUserID, auth.GetUserID(c))
		c.Next()
	}
}
