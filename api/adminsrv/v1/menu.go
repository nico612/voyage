package v1

import (
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/pkg/validator"
)

type BaseMenuReq struct {
	models.SysBaseMenu
}

func (m BaseMenuReq) Validate() error {

	rules := validator.Rules{
		"Path":      {validator.NotEmpty()},
		"ParentId":  {validator.NotEmpty()},
		"Name":      {validator.NotEmpty()},
		"Component": {validator.NotEmpty()},
		"Sort":      {validator.Ge("0")},
	}

	return validator.Verify(m.SysBaseMenu, rules)
}

type SysMenusResponse struct {
	Menus []models.SysMenu `json:"menus"`
}
