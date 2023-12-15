package models

import (
	"gorm.io/gorm"
	"time"
)

// 角色
type SysAuthority struct {
	CreatedAt       time.Time       // 创建时间
	UpdatedAt       time.Time       // 更新时间
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`                                          // 删除时间
	AuthorityId     uint            `json:"authorityId" gorm:"primaryKey;comment:角色ID;size:90"`      // 角色ID
	AuthorityName   string          `json:"authorityName" gorm:"comment:角色名"`                        // 角色名
	ParentId        *uint           `json:"parentId" gorm:"comment:父角色ID"`                           // 父角色ID
	DataAuthorityId []*SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_data_authority_id;"` // 资源授权，不建议使用该功能
	Children        []SysAuthority  `json:"children" gorm:"-"`                                       // 子角色
	SysBaseMenus    []SysBaseMenu   `json:"menus" gorm:"many2many:sys_authority_menus;"`             // 路由菜单
	Users           []SysUser       `json:"-" gorm:"many2many:sys_user_authority;"`                  // 用户
	DefaultRouter   string          `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"`     // 默认菜单(默认dashboard)

}

func (SysAuthority) TableName() string {
	return "sys_authorities"
}
