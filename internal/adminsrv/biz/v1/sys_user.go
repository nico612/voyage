package v1

import (
	"context"
	"github.com/jinzhu/copier"
	v1 "github.com/nico612/voyage/api/adminsrv/v1"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/pkg/auth"
	"github.com/nico612/voyage/pkg/errors"
)

type SysUserService interface {
	GetUserWithID(ctx context.Context, userID uint) (*models.SysUser, error)
	GetUserWithUsername(ctx context.Context, username string) (*models.SysUser, error)
	Update(ctx context.Context, user *models.SysUser) error
	Register(ctx context.Context, req *v1.UserRegisterReq) (user *models.SysUser, err error)
	SetUserAuthorities(ctx context.Context, req *v1.SetUserAuthoritiesReq) error
	SetUserInfo(ctx context.Context, req *v1.ChangeUserInfoReq) error
	GetUserInfoList(ctx context.Context, pageInfo *v1.PageInfo) ([]models.SysUser, int64, error)
	SetUserAuthority(ctx context.Context, userId uint, authorityId uint) error
	DeleteUser(ctx context.Context, userId uint) error
}

var _ SysUserService = (*sysUserService)(nil)

// 提供 User Service
func newUserService(s *service) SysUserService {
	return &sysUserService{store: s.store}
}

type sysUserService struct {
	store store.IStore
}

func (u *sysUserService) DeleteUser(ctx context.Context, userId uint) error {

	return u.store.Transaction(ctx, func(txCtx context.Context) error {

		// 1, 删除用户信息
		if err := u.store.SysUsers(txCtx).DeleteUser(txCtx, userId); err != nil {
			return err
		}

		// 2. 删除用户角色表
		return u.store.SysUsers(txCtx).DeleteUserAuthority(txCtx, userId)
	})
}

// SetUserAuthority 设置用户角色
func (u *sysUserService) SetUserAuthority(ctx context.Context, userId uint, authorityId uint) error {
	// 1. 查询用户是否有该角色
	if !u.store.SysUsers(ctx).ExistUserAuthority(ctx, userId, authorityId) {
		return errors.New("user not has this authority")
	}

	// 2. 更新用户角色
	return u.store.SysUsers(ctx).UpdateUserAuthority(ctx, userId, authorityId)
}

func (u *sysUserService) GetUserInfoList(ctx context.Context, pageInfo *v1.PageInfo) ([]models.SysUser, int64, error) {
	return u.store.SysUsers(ctx).GetUserInfoList(ctx, pageInfo.Offset, pageInfo.Limit)
}

func (u *sysUserService) SetUserInfo(ctx context.Context, req *v1.ChangeUserInfoReq) error {

	// 1. 设置权限组
	if err := u.SetUserAuthorities(ctx, &v1.SetUserAuthoritiesReq{
		ID:           req.ID,
		AuthorityIds: req.AuthorityIds,
	}); err != nil {
		return err
	}

	// 2. 更新用户信息
	userChangeInfo := &models.SysUser{}
	_ = copier.Copy(userChangeInfo, req)

	return u.store.SysUsers(ctx).UpdateUserInfo(ctx, userChangeInfo)
}

func (u *sysUserService) SetUserAuthorities(ctx context.Context, req *v1.SetUserAuthoritiesReq) error {

	exist, err := u.store.SysUsers(ctx).UserExistWithUserID(ctx, req.ID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.WithCode(code.ErrUserNotFound, "user not found")
	}

	return u.store.Transaction(ctx, func(txCtx context.Context) error {

		// 删除原user对应的角色组
		if err := u.store.SysUsers(txCtx).DeleteUserAuthority(txCtx, req.ID); err != nil {
			return err
		}

		userAuthority := make([]models.SysUserAuthority, 0, len(req.AuthorityIds))
		for _, authorityId := range req.AuthorityIds {
			userAuthority = append(userAuthority, models.SysUserAuthority{SysUserId: req.ID, SysAuthorityAuthorityId: authorityId})
		}

		// 创建新的user对应的角色组
		if err := u.store.SysUsers(txCtx).CreateUserAuthority(txCtx, userAuthority); err != nil {
			return err
		}

		// 更新用户角色id
		if err := u.store.SysUsers(txCtx).UpdateUserAuthority(txCtx, req.ID, req.AuthorityIds[0]); err != nil {
			return err
		}
		return nil
	})
}

func (u *sysUserService) Register(ctx context.Context, req *v1.UserRegisterReq) (*models.SysUser, error) {
	user := &models.SysUser{}
	_ = copier.Copy(user, req)

	isExist, err := u.store.SysUsers(ctx).UserExistWithName(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if isExist {
		return nil, errors.WithCode(code.ErrUserAlreadyExist, "user already  exist")
	}

	// 用户角色关联
	authorities := make([]models.SysAuthority, 0, len(req.AuthorityIds))
	for _, authorityId := range req.AuthorityIds {
		authorities = append(authorities, models.SysAuthority{AuthorityId: authorityId})
	}
	user.Authorities = authorities
	user.Password, _ = auth.Encrypt(req.Password)

	// 创建用户
	err = u.store.SysUsers(ctx).CreateUser(ctx, user)

	return user, err
}

func (u *sysUserService) GetUserWithID(ctx context.Context, userID uint) (*models.SysUser, error) {
	return u.store.SysUsers(ctx).GetUserWithID(ctx, userID)
}

func (u *sysUserService) GetUserWithUsername(ctx context.Context, username string) (*models.SysUser, error) {
	return u.store.SysUsers(ctx).GetUserWithUsername(ctx, username)
}

func (u *sysUserService) Update(ctx context.Context, user *models.SysUser) error {

	return u.store.SysUsers(ctx).Update(ctx, user)
}
