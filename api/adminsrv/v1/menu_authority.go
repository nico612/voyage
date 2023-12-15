package v1

import (
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/pkg/validator"
)

type AddMenuAuthorityReq struct {
	Menus       []models.SysBaseMenu `json:"menus"`       // 路由菜单
	AuthorityId uint                 `json:"authorityId"` // 角色Id
}

func (r AddMenuAuthorityReq) Validate() error {
	authorityIdVerify := validator.Rules{"AuthorityId": {validator.NotEmpty()}}
	return validator.Verify(r, authorityIdVerify)
}

// DefaultMenu 角色默认的菜单权限
func DefaultMenu() []models.SysBaseMenu {
	return []models.SysBaseMenu{{
		BaseModel: models.BaseModel{ID: 1},
		ParentId:  "0",
		Path:      "dashboard",
		Name:      "dashboard",
		Component: "view/dashboard/index.vue",
		Sort:      1,
		Meta: models.Meta{
			Title: "仪表盘",
			Icon:  "setting",
		},
	}}
}
