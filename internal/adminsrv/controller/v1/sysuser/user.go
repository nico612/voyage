package sysuser

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/nico612/voyage/api/adminsrv/v1"
	"github.com/nico612/voyage/internal/adminsrv/auth"
	"github.com/nico612/voyage/internal/pkg/code"
	"github.com/nico612/voyage/internal/pkg/response"
	"github.com/nico612/voyage/pkg/errors"
	"strconv"
)

// GetUserInfo
// @Tags      SysUser
// @Summary   获取用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=v1.SysUserResp,msg=string}  "获取用户信息"
// @Router    /user/getUserInfo [get]
func (ctr *SysUserController) GetUserInfo(c *gin.Context) {
	userID := auth.GetUserID(c)
	user, err := ctr.srv.SysUsers().GetUserWithID(c, userID)
	if err != nil {
		response.Failed(c, err)
		return
	}
	resp := v1.SysUserResp{
		SysUser: *user,
	}
	response.Success(c, resp)
}

// Register
// @Tags     SysUser
// @Summary  用户注册账号
// @Produce   application/json
// @Param    data  body      v1.UserRegisterReq                                            true  "用户名, 昵称, 密码, 角色ID"
// @Success  200   {object}  response.Response{data=v1.SysUserResp,msg=string}  "用户注册账号,返回包括用户信息"
// @Router   /user/admin_register [post]
func (ctr *SysUserController) Register(c *gin.Context) {
	var r v1.UserRegisterReq

	if err := c.ShouldBindJSON(&r); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind error: %s", err.Error()))
		return
	}

	if err := r.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validate error: %s", err.Error()))
		return
	}

	user, err := ctr.srv.SysUsers().Register(c, &r)
	if err != nil {
		response.Failed(c, err)
		return
	}
	response.Success(c, v1.SysUserResp{SysUser: *user})
}

func (ctr *SysUserController) ChangePassword(c *gin.Context) {

}

// SetUserAuthority
// @Tags      SysUser
// @Summary   更改用户权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.SetUserAuthReq          true  "用户UUID, 角色ID"
// @Success   200   {object}  response.Response{msg=string}  "设置用户权限"
// @Router    /user/setUserAuthority [post]
func (ctr *SysUserController) SetUserAuthority(c *gin.Context) {
	var req v1.SetUserAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind error: %s", err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validate error: %s", err.Error()))
		return
	}

	// 获取用户ID
	claims, err := auth.GetClaims(c)
	if err != nil {
		response.Failed(c, errors.WithCode(code.ErrTokenInvalid, "auth get claims err: %s", err.Error()))
		return
	}

	// 更改用户角色
	if err := ctr.srv.SysUsers().SetUserAuthority(c, claims.UserID, req.AuthorityId); err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "set user authority err: %s", err.Error()))
		return
	}

	// 重新签发 token
	auth := auth.NewJwtAuth([]byte(ctr.cfg.JwtOptions.Key))
	claims.AuthorityId = req.AuthorityId
	token, err := auth.GeneratorToken(claims)
	if err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "generator token error: %s", err.Error()))
		return
	}

	c.Header("new-token", token)
	c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt.Unix(), 10))
	response.Success(c)
}

// DeleteUser
// @Tags      SysUser
// @Summary   删除用户
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.GetById                true  "用户ID"
// @Success   200   {object}  response.Response{msg=string}  "删除用户"
// @Router    /user/deleteUser [delete]
func (ctr *SysUserController) DeleteUser(c *gin.Context) {
	var req v1.GetById

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind error: %s", err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validate error: %s", err.Error()))
		return
	}

	if err := ctr.srv.SysUsers().DeleteUser(c, req.ID); err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "delete user error: %s", err.Error()))
		return
	}

	response.Success(c)
}

// SetUserInfo
// @Tags      SysUser
// @Summary   设置用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.ChangeUserInfoReq                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  response.Response{data=string,msg=string}  "设置用户信息"
// @Router    /user/setUserInfo [put]
func (ctr *SysUserController) SetUserInfo(c *gin.Context) {
	var req v1.ChangeUserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind err: %s", err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		response.Failed(c, errors.WithCode(code.ErrValidation, "validation err: %s", err.Error()))
		return
	}

	if err := ctr.srv.SysUsers().SetUserInfo(c, &req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "validation err: %s", err.Error()))
		return
	}

	response.Success(c)
}

func (ctr *SysUserController) SetSelfInfo(c *gin.Context) {

}

// SetUserAuthorities
// @Tags      SysUser
// @Summary   设置用户权限组
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.SetUserAuthoritiesReq   true  "用户UUID, 角色ID"
// @Success   200   {object}  response.Response{msg=string}  "设置用户权限"
// @Router    /user/setUserAuthorities [post]
func (ctr *SysUserController) SetUserAuthorities(c *gin.Context) {
	var req v1.SetUserAuthoritiesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind json err: %s", err.Error()))
		return
	}

	if err := ctr.srv.SysUsers().SetUserAuthorities(c, &req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "set user authorities err: %s", err.Error()))
		return
	}

	response.Success(c)
}

func (ctr *SysUserController) ResetPassword(c *gin.Context) {

}

// GetUserList
// @Tags      SysUser
// @Summary   分页获取用户列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      v1.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=v1.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /user/getUserList [post]
func (ctr *SysUserController) GetUserList(c *gin.Context) {
	var req v1.PageInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, errors.WithCode(code.ErrBind, "bind json err: %s", err.Error()))
		return
	}

	_ = req.Validate()

	list, total, err := ctr.srv.SysUsers().GetUserInfoList(c, &req)
	if err != nil {
		response.Failed(c, errors.WithCode(code.ErrUnknown, "get user info list err: %s", err.Error()))
		return
	}

	response.Success(c, v1.PageResult{
		List: list,
		PageInfo: v1.PageInfo{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	})

}
