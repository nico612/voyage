package v1

import (
	"context"
	"github.com/jinzhu/copier"
	v1 "github.com/nico612/voyage/api/adminsrv/v1"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/adminsrv/store"
)

type SysBaseMenuService interface {

	// CreateBaseMenu 添加基础路菜单
	CreateBaseMenu(ctx context.Context, menu *v1.BaseMenuReq) error

	// DeleteBaseMenu 删除基础菜单
	DeleteBaseMenu(ctx context.Context, menuId uint) error

	// UpdateBaseMenu 更新基础菜单
	UpdateBaseMenu(ctx context.Context, menu *v1.BaseMenuReq) error

	// GetBaseMenuById 获取基础菜单
	GetBaseMenuById(ctx context.Context, menuId uint) (*models.SysBaseMenu, error)

	// GetBaseMenuInfoList 获取基础菜单列表
	GetBaseMenuInfoList(ctx context.Context, info *v1.PageInfo) ([]models.SysBaseMenu, error)
}

type sysBaseMenuService struct {
	store store.IStore
}

func (s *sysBaseMenuService) GetBaseMenuInfoList(ctx context.Context, info *v1.PageInfo) ([]models.SysBaseMenu, error) {
	return s.store.SysBaseMenus(ctx).GetBaseMenuTree(ctx)
}

func (s *sysBaseMenuService) CreateBaseMenu(ctx context.Context, menu *v1.BaseMenuReq) error {
	m := &models.SysBaseMenu{}
	_ = copier.Copy(m, menu)
	return s.store.SysBaseMenus(ctx).CreateBaseMenu(ctx, m)
}

func (s *sysBaseMenuService) DeleteBaseMenu(ctx context.Context, menuId uint) error {
	return s.store.SysBaseMenus(ctx).DeleteBaseMenu(ctx, menuId)
}

func (s *sysBaseMenuService) UpdateBaseMenu(ctx context.Context, menu *v1.BaseMenuReq) error {
	baseMenu := menu.SysBaseMenu
	return s.store.SysBaseMenus(ctx).UpdateBaseMenu(ctx, &baseMenu)
}

func (s *sysBaseMenuService) GetBaseMenuById(ctx context.Context, menuId uint) (*models.SysBaseMenu, error) {
	return s.store.SysBaseMenus(ctx).GetBaseMenuById(ctx, menuId)
}

func newBaseSysMenuService(store store.IStore) SysBaseMenuService {
	return &sysBaseMenuService{store: store}
}
