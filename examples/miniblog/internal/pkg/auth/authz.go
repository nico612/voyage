// Copyright 2023 Innkeeper Belm(孔令飞) <nosbelm@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package auth

import (
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// 基于权限控制列表 ACL

const (
	// casbin 访问控制模型.
	aclModel = `[request_definition]
r =  sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`
)

// Authz 定义了一个授权器，提供授权功能.
type Authz struct {
	*casbin.SyncedEnforcer
}

// NewAuthz 创建一个使用 casbin 完成授权的授权器.
func NewAuthz(db *gorm.DB) (*Authz, error) {
	// Initialize a Gorm adapter and use it in a Casbin enforcer
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	m, err := model.NewModelFromString(aclModel)
	if err != nil {
		return nil, err
	}

	// Initialize the enforcer.
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// Load the policy form DB
	if err = enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	// 为了能够快速将数据库中保存的授权策略加载到缓存中，供授权时使用，这里我们设置了每 5s 从数据库中同步一次授权策略：
	enforcer.StartAutoLoadPolicy(5 * time.Second)

	a := &Authz{
		enforcer,
	}

	return a, nil
}

// Authorize 用来进行授权.
func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	return a.Enforce(sub, obj, act)
}
