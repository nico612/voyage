package menu

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/nico612/voyage/api/adminsrv/v1"
	"github.com/nico612/voyage/internal/adminsrv/auth"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/internal/pkg/response"
	"github.com/nico612/voyage/pkg/errors"
)

// AddMenuAuthority
// @Tags      AuthorityMenu
// @Summary   增加menu和角色关联关系
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.AddMenuAuthorityReq  true  "角色ID"
// @Success   200   {object}  response.Response{msg=string}   "增加menu和角色关联关系"
// @Router    /menu/addMenuAuthority [post]
func (a *AuthorityMenuController) AddMenuAuthority(c *gin.Context) {
	var req v1.AddMenuAuthorityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind json error: %s", err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validate error: %s", err.Error()))
		return
	}

	if err := a.srv.SysAuthority().SetMenuAuthority(c, req.Menus, req.AuthorityId); err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "add menu authority err: %s", err.Error()))
		return
	}

	response.Success(c)
}

// GetMenu
// @Tags      AuthorityMenu
// @Summary   获取用户动态路由
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      v1.Empty                                                  true  "空"
// @Success   200   {object}  response.Response{data=v1.SysMenusResponse,msg=string}  "获取用户动态路由,返回包括系统菜单详情列表"
// @Router    /menu/getMenu [post]
func (a *AuthorityMenuController) GetMenu(c *gin.Context) {

	menus, err := a.srv.SysMenu().GetMenuTree(c, auth.GetUserAuthorityID(c))

	if err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "get menu tree error: %s", err.Error()))
		return
	}
	if menus == nil {
		menus = make([]models.SysMenu, 0)
	}

	response.Success(c, v1.SysMenusResponse{Menus: menus})
}

func (a *AuthorityMenuController) GetBaseMenuTree(c *gin.Context) {

}

// GetMenuAuthority
// @Tags      AuthorityMenu
// @Summary   获取指定角色menu
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetAuthorityId                                     true  "角色ID"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "获取指定角色menu"
// @Router    /menu/getMenuAuthority [post]
func (a *AuthorityMenuController) GetMenuAuthority(c *gin.Context) {

}

func (a *AuthorityMenuController) GetBaseMenuById(c *gin.Context) {

}
