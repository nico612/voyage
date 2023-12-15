package mycasbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/nico612/voyage/pkg/errors"
	"github.com/nico612/voyage/pkg/log"
	"gorm.io/gorm"
	"sync"
)

// Casbin 权限管理
var text = `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`

var (
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
	once                 sync.Once
)

func CreateCasbinOr(db *gorm.DB) *casbin.SyncedCachedEnforcer {
	if syncedCachedEnforcer == nil {
		once.Do(func() {
			apter, err := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rule")
			if err != nil {
				log.Errorf("casbin new adapter error: %s", err.Error())
				return
			}

			// 模型
			m, err := model.NewModelFromString(text)
			if err != nil {
				log.Errorf("casbin new model error: %s", err.Error())
				return
			}

			syncedCachedEnforcer, err = casbin.NewSyncedCachedEnforcer(m, apter)
			if err != nil {
				log.Errorf("casbin new synced cache enforcer error: %s", err.Error())
				return
			}
			syncedCachedEnforcer.SetExpireTime(60 * 60)
			_ = syncedCachedEnforcer.LoadPolicy()

			// TODO 使用 redis 监听策略变化然后更新策略

		})
	}

	return syncedCachedEnforcer
}

// FreshCasbin 重新加载策略， 在用户重新添加或生成的时候需要重新添加
func FreshCasbin() error {
	if syncedCachedEnforcer == nil {
		return errors.New("syncedCachedEnforcer is nil")
	}
	return syncedCachedEnforcer.LoadPolicy()
}
