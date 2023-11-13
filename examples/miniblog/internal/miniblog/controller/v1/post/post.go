// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package post

import (
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/store"
)

// PostController 是 post 模块在 Controller 层的实现，用来处理用户模块的请求.
// 如果要实现grpc 注意 PostController 必须内嵌 pb.UnimplementedMiniBlogServer 结构体，否则编译时会报错：missing mustEmbedUnimplementedMiniBlogServer method...
type PostController struct {
	b biz.IBiz // 依赖业务层
}

// New 创建一个 user controller.
func New(ds store.IStore) *PostController {
	// Controller 依赖 Biz，Biz 依赖 Store，所以我们传入了 IStore 类型的参数 ds 用来创建 PostController。
	return &PostController{b: biz.NewBiz(ds)}
}
