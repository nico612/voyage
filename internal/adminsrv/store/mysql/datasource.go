package mysql

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/pkg/errors"
	"gorm.io/gorm"
)

type txCtxKey struct{}

func TxFromContext(ctx context.Context) *gorm.DB {
	tx, _ := ctx.Value(txCtxKey{}).(*gorm.DB)
	return tx
}

func NewTxContext(parent context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(parent, txCtxKey{}, tx)
}

// 提供数据储存功能，实现 store.IStore 接口
type datasource struct {
	db *gorm.DB
}

func (ds *datasource) SysCasbin(ctx context.Context) store.SysCasbinStore {

	return newSysCasbin(ds.dbFromContext(ctx))
}

func (ds *datasource) SysAuthority(ctx context.Context) store.SysAuthorityStore {
	return newSysAuthority(ds.dbFromContext(ctx))
}

func (ds *datasource) SysBaseMenus(ctx context.Context) store.SysBaseMenuStore {
	return newSysBaseMenu(ds.dbFromContext(ctx))
}

func (ds *datasource) SysUsers(ctx context.Context) store.SysUserStore {
	return newUsers(ds.dbFromContext(ctx))
}

func (ds *datasource) JwtBlack(ctx context.Context) store.JwtBlackStore {
	return newJwtBlack(ds.dbFromContext(ctx))
}

func (ds *datasource) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}
	return db.Close()
}

// dbFromContext 从context中获取db实例
func (ds *datasource) dbFromContext(parent context.Context) *gorm.DB {
	tx := TxFromContext(parent)
	if tx != nil {
		return tx
	}
	return ds.db
}

// Transaction 开启事务
func (ds *datasource) Transaction(ctx context.Context, fc func(txCtx context.Context) error) error {
	tx := TxFromContext(ctx)
	// 已存在事务，则直接调用 fc
	if tx != nil {
		return fc(ctx)
	}
	return ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 储存到 上下文中，并执行回调函数你
		return fc(NewTxContext(ctx, tx))
	})
}
