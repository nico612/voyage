package v1

import (
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/pkg/validator"
)

type SysUserResp struct {
	models.SysUser
}

type UserRegisterReq struct {
	Username     string `json:"userName" form:"userName"`         // 用户名
	NickName     string `json:"nickName" form:"nickname"`         // 昵称
	Password     string `json:"password" form:"password"`         // 密码
	AuthorityId  uint   `json:"authorityId" form:"authorityId"`   // 角色id
	HeaderImg    string `json:"headerImg" form:"headerImg"`       // 头像
	AuthorityIds []uint `json:"authorityIds" form:"authorityIds"` // 多个角色id
	Phone        string `json:"phone" form:"phone"`               // 电话
	Email        string `json:"email" form:"email"`               // 邮箱
	Enable       int    `json:"enable"`                           // 是否启用
}

func (u UserRegisterReq) Validate() error {
	rules := validator.Rules{"Username": {validator.NotEmpty()}, "NickName": {validator.NotEmpty()}, "Password": {validator.NotEmpty()}, "AuthorityId": {validator.NotEmpty()}}
	return validator.Verify(u, rules)
}

type SetUserAuthoritiesReq struct {
	ID           uint   // 用户ID
	AuthorityIds []uint `json:"authorityIds"` // 角色ID
}

// ChangeUserInfoReq 修改用户信息
type ChangeUserInfoReq struct {
	ID           uint                  `json:"id"`           // 主键ID
	NickName     string                `json:"nickName"`     // 用户昵称
	Phone        string                `json:"phone"`        // 用户手机号
	AuthorityIds []uint                `json:"authorityIds"` // 角色ID
	Email        string                `json:"email"`        // 用户邮箱
	HeaderImg    string                `json:"headerImg"`    // 用户头像
	SideMode     string                `json:"sideMode"`     // 用户侧边主题
	Enable       int                   `json:"enable"`       //冻结用户
	Authorities  []models.SysAuthority `json:"authorities"`  // 权限组
}

func (c ChangeUserInfoReq) Validate() error {
	IdVerify := validator.Rules{"ID": []string{validator.NotEmpty()}}
	return validator.Verify(c, IdVerify)
}

// SetUserAuthReq 设置用户角色
type SetUserAuthReq struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"`
}

func (u SetUserAuthReq) Validate() error {
	SetUserAuthorityVerify := validator.Rules{"AuthorityId": {validator.NotEmpty()}}
	return validator.Verify(u, SetUserAuthorityVerify)
}
