package mysql

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/pkg/errors"
	"gorm.io/gorm"
	"strconv"
)

var (
	_ store.SysBaseMenuStore = (*sysBaseMenu)(nil)
)

func newSysBaseMenu(db *gorm.DB) *sysBaseMenu {
	return &sysBaseMenu{
		db: db,
	}
}

type sysBaseMenu struct {
	db *gorm.DB
}

func (s *sysBaseMenu) GetBaseMenuByMenuIds(ctx context.Context, ids []string) (menus []models.SysBaseMenu, err error) {
	err = s.db.Preload("Parameters").Order("sort").Where("id in (?)", ids).Find(&menus).Error
	return
}

func (s *sysBaseMenu) GetBaseMenuTree(ctx context.Context) (menus []models.SysBaseMenu, err error) {
	// 获取路由树
	treeMap, err := s.getBaseMenuTreeMap(ctx)
	if err != nil {
		return nil, err
	}
	menus = treeMap["0"] // 根路由

	// 递归获取嵌套的树
	for i := 0; i < len(menus); i++ {
		err = s.getChildrenList(ctx, &menus[i], treeMap)
	}
	return menus, err
}

// getChildrenList 获取菜单的子菜单列表
func (s *sysBaseMenu) getChildrenList(ctx context.Context, menu *models.SysBaseMenu, treeMap map[string][]models.SysBaseMenu) (err error) {
	// 根据路由树递归获取子菜单
	menu.Children = treeMap[strconv.Itoa(int(menu.ID))]
	// 递归
	for i := 0; i < len(menu.Children); i++ {
		err = s.getChildrenList(ctx, &menu.Children[i], treeMap)
	}
	return err
}

// getBaseMenuTreeMap 获取路由总树map，将基础路由按照相同的父路由进行分组
func (s *sysBaseMenu) getBaseMenuTreeMap(ctx context.Context) (treeMap map[string][]models.SysBaseMenu, err error) {
	var allMenus []models.SysBaseMenu
	// 初始化treeMap
	treeMap = make(map[string][]models.SysBaseMenu)

	if err = s.db.Order("sort DESC").Preload("MenuBtn").Preload("Parameters").Find(&allMenus).Error; err != nil {
		return
	}
	for _, menu := range allMenus {
		treeMap[menu.ParentId] = append(treeMap[menu.ParentId], menu)
	}
	return treeMap, err
}

func (s *sysBaseMenu) DeleteBaseMenu(ctx context.Context, menuId uint) error {

	// 查找该菜单是否存在子菜单
	err := s.db.Preload("MenBtn").Preload("Parameters").Where("parent_id = ?", menuId).First(&models.SysBaseMenu{}).Error
	if err != nil {
		var menu models.SysBaseMenu
		// 删除菜单
		db := s.db.Preload("SysAuthoritys").Where("id = ?", menuId).First(&menu).Delete(&menu)
		// 删除 SysBaseMenuParameter 库 中记录
		err = s.db.Delete(&models.SysBaseMenuParameter{}, "sys_base_menu_id = ?", menuId).Error
		// 删除 SysBaseMenuBtn 表记录
		err = s.db.Delete(&models.SysBaseMenuBtn{}, "sys_base_menu_id = ?", menuId).Error
		// 删除 SysAuthorityBtn 表记录
		err = s.db.Delete(&models.SysAuthorityBtn{}, "sys_menu_id = ?", menuId).Error

		if len(menu.SysAuthoritys) > 0 {
			// 删除关联记录
			err = s.db.Model(&menu).Association("SysAuthoritys").Delete(&menu.SysAuthoritys)
		} else {
			err = db.Error
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("此菜单存在子菜单，不能删除")
	}

	return err

}

func (s *sysBaseMenu) UpdateBaseMenu(ctx context.Context, menu *models.SysBaseMenu) error {
	var oldMenu models.SysBaseMenu
	upDateMap := make(map[string]interface{})
	upDateMap["keep_alive"] = menu.KeepAlive
	upDateMap["close_tab"] = menu.CloseTab
	upDateMap["default_menu"] = menu.DefaultMenu
	upDateMap["parent_id"] = menu.ParentId
	upDateMap["path"] = menu.Path
	upDateMap["name"] = menu.Name
	upDateMap["hidden"] = menu.Hidden
	upDateMap["component"] = menu.Component
	upDateMap["title"] = menu.Title
	upDateMap["active_name"] = menu.ActiveName
	upDateMap["icon"] = menu.Icon
	upDateMap["sort"] = menu.Sort
	err := s.db.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", menu.ID).Find(&oldMenu)
		// 1. 判断是否重名
		if oldMenu.Name != menu.Name {
			if txErr := tx.Where("id <> ? AND name = ?", menu.ID, menu.Name).First(&sysBaseMenu{}).Error; !errors.Is(txErr, gorm.ErrRecordNotFound) {
				return errors.New("存在相同name修改失败")
			}
		}

		// 删除原 SysBaseMenuParameter
		if txErr := tx.Unscoped().Delete(&models.SysBaseMenuParameter{}, "sys_base_menu_id = ?", menu.ID).Error; txErr != nil {
			return txErr
		}

		// 删除原 SysBaseMenuBtn
		if txErr := tx.Unscoped().Delete(&models.SysBaseMenuBtn{}, "sys_base_menu_id = ?", menu.ID).Error; txErr != nil {
			return txErr
		}

		// 新创建 SysBaseMenuParameter
		if len(menu.Parameters) > 0 {
			for i := range menu.Parameters {
				menu.Parameters[i].SysBaseMenuID = menu.ID
			}
			if txErr := tx.Create(&menu.Parameters).Error; txErr != nil {
				return txErr
			}
		}

		// 新创建 SysBaseMenuBtn
		if len(menu.MenuBtn) > 0 {
			for i := range menu.MenuBtn {
				menu.MenuBtn[i].SysBaseMenuID = menu.ID
			}
			if txErr := tx.Create(&menu.MenuBtn).Error; txErr != nil {
				return txErr
			}
		}

		// 更新
		if txErr := db.Updates(upDateMap).Error; txErr != nil {
			return txErr
		}
		return nil
	})
	return err
}

func (s *sysBaseMenu) GetBaseMenuById(ctx context.Context, menuId uint) (*models.SysBaseMenu, error) {
	menu := &models.SysBaseMenu{}
	if err := s.db.Preload("MenuBtn").Preload("Parameters").Where("id = ?", menuId).First(&menu).Error; err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *sysBaseMenu) CreateBaseMenu(ctx context.Context, menu *models.SysBaseMenu) error {
	// 查找并创建，如果有相同的值，只做查询操作，不会报错，如果没有则创建
	// s.db.Where("name = ?", menu.Name).FirstOrCreate(&menu).Error
	if err := s.db.Where("name = ?", menu.Name).First(&models.SysBaseMenu{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("存在相同的name")
	}

	return s.db.Create(&menu).Error
}
