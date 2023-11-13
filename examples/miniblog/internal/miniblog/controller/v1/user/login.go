// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/core"
	v1 "github.com/nico612/go-project/examples/miniblog/pkg/api/miniblog/v1"
)

func (ctrl *UserController) Login(c *gin.Context) {
	var r v1.LoginRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	resp, err := ctrl.b.Users().Login(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}
