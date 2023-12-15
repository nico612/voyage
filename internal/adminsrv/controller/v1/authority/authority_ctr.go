package authority

import (
	"github.com/gin-gonic/gin"
	"github.com/nico612/voyage/api/adminsrv/v1"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/internal/pkg/response"
	"github.com/nico612/voyage/pkg/errors"
)

// CreateAuthority
// @Tags      Authority
// @Summary   创建角色
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.SysCreateOrUpdateAuthorityReq  true  "权限id, 权限名, 父角色id"
// @Success   200   {object}  response.Response{data=v1.SysAuthorityResponse,msg=string}  "创建角色,返回包括系统角色详情"
// @Router    /authority/createAuthority [post]
func (ctr *AuthorityController) CreateAuthority(c *gin.Context) {
	var req v1.SysCreateOrUpdateAuthorityReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind error: %s", err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validate error: %s", err.Error()))
	}

	authority, err := ctr.srv.SysAuthority().CreateAuthority(c, &req)
	if err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "create authority error: %s", err.Error()))
		return
	}

	response.Success(c, authority)

}

func (ctr *AuthorityController) DeleteAuthority(c *gin.Context) {

}

// UpdateAuthority
// @Tags      Authority
// @Summary   更新角色信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.SysCreateOrUpdateAuthorityReq                                             true  "权限id, 权限名, 父角色id"
// @Success   200   {object}  response.Response{data=v1.Empty, msg=string}  "更新角色信息,返回包括系统角色详情"
// @Router    /authority/updateAuthority [post]
func (ctr *AuthorityController) UpdateAuthority(c *gin.Context) {
	var req v1.SysCreateOrUpdateAuthorityReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind error: %s", err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validate error: %s", err.Error()))
		return
	}

	if err := ctr.srv.SysAuthority().UpdateAuthority(c, &req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "update authority errr: %s", err.Error()))
		return
	}

	response.Success(c, v1.Empty{})
}

func (ctr *AuthorityController) CopyAuthority(c *gin.Context) {

}

func (ctr *AuthorityController) SetDataAuthority(c *gin.Context) {

}

// GetAuthorityList
// @Tags      Authority
// @Summary   分页获取角色列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=v1.PageResult,msg=string}  "分页获取角色列表,返回包括列表,总数,页码,每页数量"
// @Router    /authority/getAuthorityList [post]
func (ctr *AuthorityController) GetAuthorityList(c *gin.Context) {
	var req v1.PageInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind error: %s", err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validate error: %s", err.Error()))
		return
	}

	list, total, err := ctr.srv.SysAuthority().GetAuthorityInfoList(c, &req)
	if err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "get authority error: %s", err.Error()))
		return
	}

	resPageInfo := req
	resPageInfo.Total = total
	response.Success(c, v1.PageResult{
		List:     list,
		PageInfo: resPageInfo,
	})
}
