package user

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/store"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/auth"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/errno"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/log"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/model"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/token"
	v1 "github.com/nico612/go-project/examples/miniblog/pkg/api/miniblog/v1"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"regexp"
	"sync"
)

// UserBiz 的创建思路和 UserStore 保持一致。在 Create 方法中，实现了具体的创建逻辑：
// 接受来自 Controller 层的入参：context.Context、*v1.CreateUserRequest；
//
//go:generate mockgen -destination mock_user.go -package user github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz/user UserBiz
type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error
	Get(ctx context.Context, username string) (*v1.GetUserResponse, error)
	List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error)
	Update(ctx context.Context, username string, r *v1.UpdateUserRequest) error
	Delete(ctx context.Context, username string) error
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

func (u *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {

	// 获取登录用户的所有信息
	user, err := u.ds.Users().Get(ctx, r.Username)
	if err != nil {
		return nil, err
	}

	// 对比传入的明文密码和数据库中已加密过的密码是否匹配
	if err = auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}

	// 如果匹配成功，说明登录成功，签发 token 并返回
	t, err := token.Sign(r.Username)
	if err != nil {
		return nil, errno.ErrSignToken
	}

	return &v1.LoginResponse{Token: t}, nil
}

func (u *userBiz) ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error {
	// 获取登录用户的所有信息
	user, err := u.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}

	if err = auth.Compare(user.Password, r.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}

	user.Password, _ = auth.Encrypt(r.NewPassword)
	if err = u.ds.Users().Update(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *userBiz) Get(ctx context.Context, username string) (*v1.GetUserResponse, error) {
	user, err := u.ds.Users().Get(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}
		return nil, err
	}

	var resp v1.GetUserResponse
	_ = copier.Copy(&resp, user)

	resp.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
	resp.UpdatedAt = user.UpdatedAt.Format("2006-01-02 15:04:05")
	return &resp, nil
}

func (u *userBiz) List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error) {
	count, list, err := u.ds.Users().List(ctx, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list users from storage", "err", err)
		return nil, err
	}

	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)

	// 使用 goroutine 提高接口性能
	for _, item := range list {
		user := item
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				postCount, err := u.ds.Posts().Count(ctx, user.Username)
				if err != nil {
					log.C(ctx).Errorw("Failed to list posts", "err", err)
					return err
				}
				m.Store(user.ID, &v1.UserInfo{
					Username:  user.Username,
					Nickname:  user.Nickname,
					Email:     user.Email,
					Phone:     user.Phone,
					PostCount: postCount,
					CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
					UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
				})
				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.C(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}

	users := make([]*v1.UserInfo, 0, len(list))
	for _, item := range list {
		user, _ := m.Load(item.ID)
		users = append(users, user.(*v1.UserInfo))
	}

	log.C(ctx).Debugw("Get users form backend storage", "count", len(users))

	return &v1.ListUserResponse{Users: users, TotalCount: count}, nil
}

// ListWithBadPerformance 是一个性能较差的实现方式（已废弃）.
func (b *userBiz) ListWithBadPerformance(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error) {
	count, list, err := b.ds.Users().List(ctx, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list users from storage", "err", err)
		return nil, err
	}

	users := make([]*v1.UserInfo, 0, len(list))
	for _, item := range list {
		user := item

		postCount, err := b.ds.Posts().Count(ctx, user.Username)
		if err != nil {
			log.C(ctx).Errorw("Failed to list posts", "err", err)
			return nil, err
		}

		users = append(users, &v1.UserInfo{
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Email,
			PostCount: postCount,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	log.C(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &v1.ListUserResponse{TotalCount: count, Users: users}, nil
}

func (u *userBiz) Update(ctx context.Context, username string, r *v1.UpdateUserRequest) error {
	user, err := u.ds.Users().Get(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ErrUserNotFound
		}
		return err
	}

	if r.Nickname != nil {
		user.Nickname = *r.Nickname
	}
	if r.Email != nil {
		user.Email = *r.Email
	}

	if r.Phone != nil {
		user.Phone = *r.Phone
	}

	if err = u.ds.Users().Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *userBiz) Delete(ctx context.Context, username string) error {
	if err := u.ds.Users().Delete(ctx, username); err != nil {
		return err
	}
	return nil
}
