package store

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/models"
)

type SysUserStore interface {
	// GetUserWithID 根据用户ID获取用户
	GetUserWithID(ctx context.Context, userID uint) (*models.SysUser, error)
	// GetUserWithUsername 根据用户名获取用户
	GetUserWithUsername(ctx context.Context, username string) (*models.SysUser, error)
	// Update 更新用户，所有信息
	Update(ctx context.Context, user *models.SysUser) error
	// CreateUser 创建用户
	CreateUser(ctx context.Context, user *models.SysUser) error
	// UserExistWithName 根据用户名查询用户是否存在
	UserExistWithName(ctx context.Context, username string) (bool, error)
	// UserExistWithUserID 根据 ID 查询用户是否存在
	UserExistWithUserID(ctx context.Context, userId uint) (bool, error)
	// DeleteUserAuthority 删除用户角色
	DeleteUserAuthority(ctx context.Context, userId uint) error
	// CreateUserAuthority 创建用户角色组
	CreateUserAuthority(ctx context.Context, authority []models.SysUserAuthority) error
	// UpdateUserAuthority 更新用户角色
	UpdateUserAuthority(ctx context.Context, id uint, authorityId uint) error
	// UpdateUserInfo 更新用户信息
	UpdateUserInfo(ctx context.Context, info *models.SysUser) error
	// GetUserInfoList 分页获取用户列表
	GetUserInfoList(ctx context.Context, offset int, limit int) (list []models.SysUser, total int64, err error)
	// ExistUserAuthority 查询用户是否有该角色
	ExistUserAuthority(ctx context.Context, userId uint, authorityId uint) bool
	// DeleteUser 删除用户吧
	DeleteUser(ctx context.Context, userId uint) error
}
