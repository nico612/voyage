// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package biz

import (
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz/post"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz/user"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/store"
)

//go:generate mockgen -destination mock_biz.go -package biz github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz IBiz

type IBiz interface {
	Users() user.UserBiz
	Posts() post.PostBiz
}

// biz 是IBiz的一个具体实现，其依赖 store 层.
type biz struct {
	ds store.IStore
}

// 确保 biz 实现了IBiz 接口.
var _ IBiz = (*biz)(nil)

// NewBiz 创建一个 IBiz 类型的实例.
func NewBiz(ds store.IStore) *biz {
	return &biz{ds: ds}
}

// Users 返回一个实现了 UserBiz 接口的实例.
func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}

// Posts 返回一个实现了 PostBiz 接口实例.
func (b *biz) Posts() post.PostBiz {
	return post.New(b.ds)
}
