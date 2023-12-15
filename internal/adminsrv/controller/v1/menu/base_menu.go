package menu

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/nico612/voyage/api/adminsrv/v1"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/internal/pkg/response"
	"github.com/nico612/voyage/pkg/errors"
)

// GetMenuList
// @Tags      Menu
// @Summary   分页获取基础menu列表
// @Security  ApiKeyAuth BearerTokenAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=v1.PageResult,msg=string}  "分页获取基础menu列表,返回包括列表,总数,页码,每页数量"
// @Router    /menu/getMenuList [post]
func (a *AuthorityMenuController) GetMenuList(c *gin.Context) {
	var pageInfo v1.PageInfo
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind error: %s", err.Error()))
		return
	}

	if err := pageInfo.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validation error: %s", err.Error()))
		return
	}

	// 查询 list
	meunList, err := a.srv.SysBaseMenu().GetBaseMenuInfoList(c, &pageInfo)
	if err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "get base menuInfoList error: %s", err.Error()))
		return
	}

	result := v1.PageResult{
		List:     meunList,
		PageInfo: pageInfo,
	}

	response.Success(c, result)

}

// AddBaseMenu
// @Tags      Menu
// @Summary   新增菜单
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.BaseMenuReq             true  "路由path, 父菜单ID, 路由name, 对应前端文件路径, 排序标记"
// @Success   200   {object}  response.Response{msg=string}  "新增菜单"
// @Router    /menu/addBaseMenu [post]
func (a *AuthorityMenuController) AddBaseMenu(c *gin.Context) {

	var menuReq v1.BaseMenuReq
	if err := c.ShouldBindJSON(&menuReq); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind json error: %s", err.Error()))
		return
	}

	if err := menuReq.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validate error: %s", err.Error()))
		return
	}

	if err := a.srv.SysBaseMenu().CreateBaseMenu(c, &menuReq); err != nil {
		response.Failed(c, errors.WithCode(code.ErrCreateBaseMenu, "create base menu error: %s", err.Error()))
		return
	}

	response.Success(c, struct{}{})

}

// DeleteBaseMenu
// @Tags      Menu
// @Summary   删除菜单
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.GetById                true  "菜单id"
// @Success   200   {object}  response.Response{msg=string}  "删除菜单"
// @Router    /menu/deleteBaseMenu [post]
func (a *AuthorityMenuController) DeleteBaseMenu(c *gin.Context) {
	var req v1.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind json error: %s", err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validate error: %s", err.Error()))
		return
	}

	if err := a.srv.SysBaseMenu().DeleteBaseMenu(c, req.ID); err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "delete base menu: %s", err.Error()))
		return
	}
	response.Success(c)
}

// UpdateBaseMenu
// @Tags      Menu
// @Summary   更新菜单
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.BaseMenuReq             true  "路由path, 父菜单ID, 路由name, 对应前端文件路径, 排序标记"
// @Success   200   {object}  response.Response{msg=string}  "更新菜单"
// @Router    /menu/updateBaseMenu [post]
func (a *AuthorityMenuController) UpdateBaseMenu(c *gin.Context) {
	var menuReq v1.BaseMenuReq
	if err := c.ShouldBindJSON(&menuReq); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind json error: %s", err.Error()))
		return
	}

	if err := menuReq.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validate error: %s", err.Error()))
		return
	}

	if err := a.srv.SysBaseMenu().UpdateBaseMenu(c, &menuReq); err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "update base menu error: %s", err.Error()))
		return
	}
	response.Success(c)

}
