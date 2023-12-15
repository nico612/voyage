package authority

import (
	"github.com/nico612/voyage/internal/adminsrv/biz/v1"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/redis/go-redis/v9"
)

type AuthorityController struct {
	srv v1.Service
	rdb *redis.Client
}

func NewAuthorityController(store store.IStore, rdb *redis.Client) *AuthorityController {
	return &AuthorityController{
		srv: v1.NewService(store, rdb),
		rdb: rdb,
	}
}
