package mysql

import (
	"context"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	v1 "github.com/nico612/voyage/api/adminsrv/v1"
	mycasbin "github.com/nico612/voyage/internal/adminsrv/casbin"
	"gorm.io/gorm"
	"strconv"
)

func newSysCasbin(db *gorm.DB) *sysCasbin {
	return &sysCasbin{
		db:               db,
		syncCacheEnforce: mycasbin.CreateCasbinOr(db),
	}
}

type sysCasbin struct {
	db               *gorm.DB
	syncCacheEnforce *casbin.SyncedCachedEnforcer
}

func (s *sysCasbin) AddPolicies(ctx context.Context, rules [][]string) error {
	casbinRules := make([]gormadapter.CasbinRule, 0, len(rules))

	for i := range rules {
		casbinRules = append(casbinRules, gormadapter.CasbinRule{
			Ptype: "p",
			V0:    rules[i][0], // 角色 id
			V1:    rules[i][1], // api
			V2:    rules[i][2], // 请求方法
		})
	}

	return s.db.Create(&casbinRules).Error
}

// RemoveFilteredPolicy 根据角色 id 移除 策略,  需要调用FreshCasbin方法才可以在系统中即刻生效
func (s *sysCasbin) RemoveFilteredPolicy(ctx context.Context, authorityId string) error {
	return s.db.Delete(&gormadapter.CasbinRule{}, "v0 = ?", authorityId).Error
}

// UpdateCasbinApi API更新随动
func (s *sysCasbin) UpdateCasbinApi(ctx context.Context, oldPath, newPath, oldMethod, newMethod string) error {
	updateMap := map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}
	return s.db.Model(&gormadapter.CasbinRule{}).Where("v1 = ? && v2 = ?", oldPath, oldMethod).Updates(updateMap).Error
}

// GetPolicyPathByAuthorityId 获取权限列表
func (s *sysCasbin) GetPolicyPathByAuthorityId(AuthorityID uint) (pathMaps []v1.CasbinInfo) {

	authorityId := strconv.Itoa(int(AuthorityID))
	list := s.syncCacheEnforce.GetFilteredPolicy(0, authorityId)
	for _, v := range list {
		pathMaps = append(pathMaps, v1.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

// FreshCasbin 刷新策略, 只要修改了策略就要应该调用该方法, 刷新策略
func (s *sysCasbin) FreshCasbin() error {
	return s.syncCacheEnforce.LoadPolicy()
}

// ClearCasbin 清除匹配的策略
func (s *sysCasbin) ClearCasbin(v int, p ...string) bool {
	success, _ := s.syncCacheEnforce.RemoveFilteredPolicy(v, p...)
	return success
}

//
//func (s *sysCasbin) UpdateCasbin(AuthorityID uint, casbinInfos []request.CasbinInfo) error {
//
//}
