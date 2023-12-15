package v1

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/pkg/log"
	"strconv"
)

// SysMenuService 动态路由菜单服务
type SysMenuService interface {
	GetMenuTree(ctx context.Context, authorityId uint) (menus []models.SysMenu, err error)
}

type sysMenuService struct {
	store store.IStore
}

func newSysMenuService(store store.IStore) SysMenuService {
	return &sysMenuService{store: store}
}

var _ SysMenuService = (*sysMenuService)(nil)

func (s *sysMenuService) getMenuTreeMap(ctx context.Context, authorityId uint) (treeMap map[string][]models.SysMenu, err error) {

	var allMenus []models.SysMenu

	var btns []models.SysAuthorityBtn

	treeMap = make(map[string][]models.SysMenu)

	// 根据id查询出所有的角色菜单
	sysAuthorityMenus, err := s.store.SysAuthority(ctx).GetSysAuthorityMenus(ctx, authorityId)
	if err != nil {
		log.L(ctx).Errorf("get sys authority  menus error: %s", err.Error())
		return
	}

	menuIds := make([]string, 0, len(sysAuthorityMenus))
	for _, menu := range sysAuthorityMenus {
		menuIds = append(menuIds, menu.MenuId)
	}

	// 根据id查询出所有的路由
	baseMenu, err := s.store.SysBaseMenus(ctx).GetBaseMenuByMenuIds(ctx, menuIds)
	if err != nil {
		log.L(ctx).Errorf("get sys base  menus error: %s", err.Error())
		return
	}

	for _, menu := range baseMenu {
		allMenus = append(allMenus, models.SysMenu{
			SysBaseMenu: menu,
			MenuId:      strconv.Itoa(int(menu.ID)),
			AuthorityId: authorityId,
			Parameters:  menu.Parameters,
		})
	}

	btns, err = s.store.SysAuthority(ctx).GetSysAuthorityBtns(ctx, authorityId)
	if err != nil {
		log.L(ctx).Errorf("get sys authority  btns error: %s", err.Error())
		return
	}

	var btnMap = make(map[uint]map[string]uint)
	for _, v := range btns {
		if btnMap[v.SysMenuID] == nil {
			btnMap[v.SysMenuID] = make(map[string]uint)
		}
		btnMap[v.SysMenuID][v.SysBaseMenuBtn.Name] = authorityId
	}

	for _, v := range allMenus {
		v.Btns = btnMap[v.ID]
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}

	return treeMap, nil
}

func (s *sysMenuService) getChildrenList(menu *models.SysMenu, treeMap map[string][]models.SysMenu) (err error) {
	menu.Children = treeMap[menu.MenuId]
	for i := 0; i < len(menu.Children); i++ {
		err = s.getChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

func (s *sysMenuService) GetMenuTree(ctx context.Context, authorityId uint) (menus []models.SysMenu, err error) {

	menuTree, err := s.getMenuTreeMap(ctx, authorityId)
	menus = menuTree["0"]

	for i := 0; i < len(menus); i++ {
		err = s.getChildrenList(&menus[i], menuTree)
	}
	return menus, err
}
