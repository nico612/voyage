package menu

import (
	v1 "github.com/nico612/voyage/internal/adminsrv/biz/v1"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/redis/go-redis/v9"
)

// AuthorityMenuController 菜单管理
type AuthorityMenuController struct {
	srv v1.Service
	rdb *redis.Client
}

func NewAuthorityMenuController(store store.IStore, rdb *redis.Client) *AuthorityMenuController {
	return &AuthorityMenuController{
		srv: v1.NewService(store, rdb),
		rdb: rdb,
	}
}
