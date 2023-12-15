package v1

import (
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/pkg/validator"
)

// 角色相关

// SysCreateOrUpdateAuthorityReq 创建角色请求
type SysCreateOrUpdateAuthorityReq struct {
	AuthorityId   uint   `json:"authorityId" form:"authorityId"`     // 角色ID
	AuthorityName string `json:"authorityName" form:"authorityName"` // 角色名
	ParentId      uint   `json:"parentId" form:"parentId"`           // 父角色Id
}

func (a SysCreateOrUpdateAuthorityReq) Validate() error {
	authorityVerify := validator.Rules{"AuthorityId": {validator.NotEmpty()}, "AuthorityName": {validator.NotEmpty()}}
	return validator.Verify(a, authorityVerify)
}

// SysAuthorityResponse 角色响应
type SysAuthorityResponse struct {
	Authority models.SysAuthority `json:"authority"`
}

// SysAuthorityCopyResponse 拷贝角色
type SysAuthorityCopyResponse struct {
	Authority      models.SysAuthority `json:"authority"`
	OldAuthorityId uint                `json:"oldAuthorityId"` // 旧角色ID
}

// SysAuthorityDeleteReq 删除
type SysAuthorityDeleteReq struct {
	AuthorityId uint `json:"authorityId"`
}

type SysAuthorityInfoReq struct {
	AuthorityId uint `json:"authorityId"`
}
