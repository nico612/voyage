package cache

import (
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/songzhibin97/gkit/options"
	"sync"
)

type LocalCache struct {
	*local_cache.Cache
}

var (
	cacheIns *LocalCache
	once     *sync.Once
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
