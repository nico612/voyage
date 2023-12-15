package store

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/models"
)

// SysBaseMenuStore 基础路由
type SysBaseMenuStore interface {

	// CreateBaseMenu 添加基础路由
	CreateBaseMenu(ctx context.Context, menu *models.SysBaseMenu) error

	// DeleteBaseMenu 删除基础路由
	DeleteBaseMenu(ctx context.Context, menuId uint) error

	// UpdateBaseMenu 更新基础路由
	UpdateBaseMenu(ctx context.Context, menu *models.SysBaseMenu) error

	// GetBaseMenuById 获取基础路由
	GetBaseMenuById(ctx context.Context, menuId uint) (*models.SysBaseMenu, error)

	// GetBaseMenuTree 获取路由列表
	GetBaseMenuTree(ctx context.Context) (menus []models.SysBaseMenu, err error)

	GetBaseMenuByMenuIds(ctx context.Context, ids []string) (menus []models.SysBaseMenu, err error)
}
