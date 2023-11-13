// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package user

import (
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/store"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/auth"
	pb "github.com/nico612/go-project/examples/miniblog/pkg/proto/miniblog/v1"
)

// UserController 是 user 模块在 Controller 层的实现，用来处理用户模块的请求.
// 在实现grpc 注意 UserController 必须内嵌 pb.UnimplementedMiniBlogServer 结构体，否则编译时会报错：missing mustEmbedUnimplementedMiniBlogServer method...
type UserController struct {
	b biz.IBiz // 依赖业务层
	a *auth.Authz
	pb.UnimplementedMiniBlogServer
}

// New 创建一个 user controller.
func New(ds store.IStore, a *auth.Authz) *UserController {
	// Controller 依赖 Biz，Biz 依赖 Store，所以我们传入了 IStore 类型的参数 ds 用来创建 UserController。
	return &UserController{b: biz.NewBiz(ds), a: a}
}
