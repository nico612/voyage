package v1

import (
	"context"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/internal/adminsrv/store/cache"
	redisStore "github.com/nico612/voyage/internal/adminsrv/store/redis"
	"github.com/redis/go-redis/v9"
	"time"
)

type JwtBlackService interface {
	// JoinJwtBlack jwt token 加入黑名单
	JoinJwtBlack(ctx context.Context, jwt *models.JwtBlackList) error

	// JwtIsInBlack 判断 jwt token 是否在黑名单中
	JwtIsInBlack(ctx context.Context, jwt string) bool

	// GetRedisJWT redis 中获取 jwt token
	GetRedisJWT(ctx context.Context, userID uint) (string, error)

	// SetRedisJWT 将 jwt token 缓存进 redis
	SetRedisJWT(ctx context.Context, userID uint, jwt string, expire time.Duration) error
}

var _ JwtBlackService = (*jwtBlackService)(nil)

type jwtBlackService struct {
	store      store.IStore
	localCache *cache.LocalCache
	rdb        *redis.Client
}

func newJwtBlackService(s *service) *jwtBlackService {
	return &jwtBlackService{
		store:      s.store,
		localCache: cache.GetLocalCacheIns(),
		rdb:        s.rdb,
	}
}

func (j *jwtBlackService) JwtIsInBlack(ctx context.Context, jwt string) bool {
	_, ok := j.localCache.Get(jwt)
	return ok
}

func (j *jwtBlackService) GetRedisJWT(ctx context.Context, userID uint) (string, error) {
	key := redisStore.GetMultiLoginLimitKey(userID)
	return j.rdb.Get(ctx, key).Result()
}

func (j *jwtBlackService) SetRedisJWT(ctx context.Context, userID uint, jwt string, expire time.Duration) error {
	key := redisStore.GetMultiLoginLimitKey(userID)
	return j.rdb.Set(ctx, key, jwt, expire).Err()
}

func (j *jwtBlackService) JoinJwtBlack(ctx context.Context, jwt *models.JwtBlackList) error {

	// 储存到数据库
	if err := j.store.JwtBlack(ctx).JoinJwtBlack(ctx, jwt); err != nil {
		return err
	}

	// 储存到 local_cache, 默认值过期时间为0，意味着直接设置为过期
	j.localCache.SetDefault(jwt.Jwt, struct{}{})

	return nil
}
