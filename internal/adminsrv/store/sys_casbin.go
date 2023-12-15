package store

import (
	"context"
	v1 "github.com/nico612/voyage/api/adminsrv/v1"
)

type SysCasbinStore interface {

	// AddPolicies 添加 cashbin 策略
	AddPolicies(ctx context.Context, rules [][]string) error

	// RemoveFilteredPolicy 移除策略使 此方法需要调用FreshCasbin方法才可以在系统中即刻生效
	RemoveFilteredPolicy(ctx context.Context, authorityId string) error

	// UpdateCasbinApi API更新随动
	UpdateCasbinApi(ctx context.Context, oldPath, newPath, oldMethod, newMethod string) error

	// GetPolicyPathByAuthorityId 获取权限列表
	GetPolicyPathByAuthorityId(AuthorityID uint) (pathMaps []v1.CasbinInfo)

	// FreshCasbin 刷新策略, 只要修改了策略就要应该调用该方法, 刷新策略
	FreshCasbin() error

	// ClearCasbin 清除匹配的策略
	ClearCasbin(v int, p ...string) bool
}
