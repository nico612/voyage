package cache

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/songzhibin97/gkit/options"
	"sync"
)

type LocalCache struct {
	*local_cache.Cache
}

var (
	cacheIns *LocalCache
	once     sync.Once
)

// GetLocalCacheIns 创建缓存实例
func GetLocalCacheIns(options ...options.Option) *LocalCache {
	if cacheIns == nil {
		once.Do(func() {
			cache := local_cache.NewCache(options...) // 获取 Cache 实例
			cacheIns = &LocalCache{
				Cache: &cache,
			}
		})
	}

	return cacheIns
}

func (l *LocalCache) LoadAllJwtBlackList(store store.IStore) error {

	data, err := store.JwtBlack(context.Background()).LoadAllJwtBlackList(context.Background())
	if err != nil {
		return err
	}
	for _, token := range data {
		l.SetDefault(token, struct{}{})
	}

	return nil
}
