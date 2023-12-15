package store

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/models"
)

type SysAuthorityStore interface {

	// ExistsAuthority 判断角色是否存在
	ExistsAuthority(ctx context.Context, authorityId uint) bool

	// GetAuthorityInfoWithId 根据 角色 id 返回角色信息
	GetAuthorityInfoWithId(ctx context.Context, authorityId uint) (*models.SysAuthority, error)

	// CreateAuthority 创建角色，同时会创建默认的菜单路由
	CreateAuthority(ctx context.Context, authority *models.SysAuthority) error

	// UpdateAuthority 更新角色
	UpdateAuthority(ctx context.Context, authority *models.SysAuthority) error

	// GetAuthorityInfoList 获取角色列表
	GetAuthorityInfoList(ctx context.Context, offset int, limit int) (list []models.SysAuthority, total int64, err error)

	// GetChildrenAuthority 获取子角色
	GetChildrenAuthority(ctx context.Context, parentId uint) (list []models.SysAuthority, err error)

	// UpdateSysBaseMenus 更新角色菜单
	UpdateSysBaseMenus(ctx context.Context, authority *models.SysAuthority) error

	// GetSysAuthorityMenus 获取角色菜单
	GetSysAuthorityMenus(ctx context.Context, authorityId uint) (authorityMenus []models.SysAuthorityMenu, err error)

	GetSysAuthorityBtns(ctx context.Context, authorityId uint) (authorityBtns []models.SysAuthorityBtn, err error)
}
