// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/core"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/errno"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/known"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/log"
)

type Auther interface {
	// Authorize 用来进行授权. sub：对象 obj：路径 act：请求方法
	Authorize(sub, obj, act string) (bool, error)
}

// Authz 权限控制列表（ACL，Access Control List）；中间件.
func Authz(a Auther) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 *gin.Context 中解析出了：用户名、访问路径、HTTP 方法，分别作为 casbin 授权模型中的 sub、obj、act。
		sub := c.GetString(known.XUsernameKey)
		obj := c.Request.URL.Path
		act := c.Request.Method
		log.Debugw("Build authorize context", "sub", sub, "obj", obj, "act", act)

		// 调用 Auther 接口的 Authorize 方法进行访问授权，授权失败返回 errno.ErrUnauthorized 错误码，并调用 c.Abort() 终止请求。
		if allowed, _ := a.Authorize(sub, obj, act); !allowed {
			core.WriteResponse(c, errno.ErrUnauthorized, nil)
			c.Abort()
			return
		}
	}
}
