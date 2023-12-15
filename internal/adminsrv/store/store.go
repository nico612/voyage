package store

import (
	"context"
)

var store IStore

type IStore interface {
	// SysUsers 后管用户接口
	SysUsers(ctx context.Context) SysUserStore
	JwtBlack(ctx context.Context) JwtBlackStore
	SysBaseMenus(ctx context.Context) SysBaseMenuStore
	SysAuthority(ctx context.Context) SysAuthorityStore
	SysCasbin(ctx context.Context) SysCasbinStore

	// Transaction 开启事务
	Transaction(ctx context.Context, fc func(txCtx context.Context) error) error
	// Close 关闭数据库连接
	Close() error
}

// Client 提供包级别的 store
func Client() IStore {
	return store
}

func SetClient(s IStore) {
	store = s
}
