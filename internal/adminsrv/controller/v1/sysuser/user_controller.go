package sysuser

import (
	v1 "github.com/nico612/voyage/internal/adminsrv/biz/v1"
	"github.com/nico612/voyage/internal/adminsrv/config"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/redis/go-redis/v9"
)

type SysUserController struct {
	store store.IStore
	cfg   *config.Config
	srv   v1.Service
	rdb   *redis.Client
}

func NewSysUserController(store store.IStore, cfg *config.Config, rdb *redis.Client) *SysUserController {
	return &SysUserController{
		store: store,
		cfg:   cfg,
		srv:   v1.NewService(store, rdb),
	}
}
