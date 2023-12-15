package redis

import (
	"context"
	"fmt"
	"github.com/nico612/voyage/internal/pkg/options"
	"github.com/nico612/voyage/pkg/errors"
	"github.com/nico612/voyage/pkg/log"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	rdb  *redis.Client
	once sync.Once

	// for lock
	lockPrefix        = "adminsrv"
	lockDefaultExpire = time.Second * 30
)

func CreateRedisOr(opts *options.RedisOptions) (*redis.Client, error) {
	if opts == nil && rdb == nil {
		return nil, errors.New("create redis error")
	}

	if rdb == nil {
		once.Do(func() {
			rdb = redis.NewClient(&redis.Options{
				Addr:         opts.Addr,
				Password:     opts.Password,
				ReadTimeout:  opts.ReadTimeout,
				WriteTimeout: opts.WriteTimeout,
				DialTimeout:  opts.DialTimeout,
				DB:           opts.DB,
				PoolSize:     10,
			})

			timeout, cancle := context.WithTimeout(context.Background(), time.Second*2)
			defer cancle()

			if err := rdb.Ping(timeout).Err(); err != nil {
				log.Fatalf("redis connect error : %s", err.Error())
			}

			log.Debug("redis connect success")
		})
	}

	return rdb, nil
}

// GetMultiLoginLimitKey 多点登录拦截key
func GetMultiLoginLimitKey(userID uint) string {
	return fmt.Sprintf("%s_%d", lockPrefix, userID)
}

//func GetLockKey(mid uint64, uid uint64) string {
//	return fmt.Sprintf("%s:%d_%d", lockPrefix, mid, uid)
//}

func Sub(channels []string) *redis.PubSub {
	return rdb.Subscribe(context.TODO(), channels...)
}
