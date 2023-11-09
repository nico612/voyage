package user

import (
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/store"
)

// UserController 是 user 模块在 Controller 层的实现，用来处理用户模块的请求.
type UserController struct {
	b biz.IBiz // 依赖业务层
}

// New 创建一个 user controller.
func New(ds store.IStore) *UserController {
	//Controller 依赖 Biz，Biz 依赖 Store，所以我们传入了 IStore 类型的参数 ds 用来创建 UserController。
	return &UserController{b: biz.NewBiz(ds)}
}
