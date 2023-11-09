package user

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/store"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/errno"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/model"
	v1 "github.com/nico612/go-project/examples/miniblog/pkg/api/miniblog/v1"
	"regexp"
)

// UserBiz 的创建思路和 UserStore 保持一致。在 Create 方法中，实现了具体的创建逻辑：
// 接受来自 Controller 层的入参：context.Context、*v1.CreateUserRequest；
type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
}

// UserBiz 接口的具体实现
type userBiz struct {
	ds store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

func (u *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)

	// 调用 store 层创建 user
	if err := u.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}
		return err
	}
	return nil
}
