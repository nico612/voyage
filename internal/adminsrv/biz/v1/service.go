package v1

import (
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/redis/go-redis/v9"
)

// Service 提供所有的服务
type Service interface {

	// SysUsers 用户管理
	SysUsers() SysUserService

	// JwtBlack jwt 黑名单管理服务
	JwtBlack() JwtBlackService

	// SysBaseMenu 基础路由管理服务
	SysBaseMenu() SysBaseMenuService

	// SysMenuAuthority 角色关联菜单
	SysMenuAuthority() SysMenuAuthorityService

	// SysMenu 动态路由管理服务
	SysMenu() SysMenuService

	// SysAuthority 角色管理服务
	SysAuthority() SysAuthorityService

	// ...
}

type service struct {
	store store.IStore
	rdb   *redis.Client
}

func (s *service) SysMenu() SysMenuService {
	return newSysMenuService(s.store)
}

func (s *service) SysAuthority() SysAuthorityService {
	return newSysAuthorityService(s.store)
}

func (s *service) SysBaseMenu() SysBaseMenuService {
	return newBaseSysMenuService(s.store)
}

func (s *service) SysMenuAuthority() SysMenuAuthorityService {
	return newSysMenuAuthorityService(s.store)
}

func (s *service) JwtBlack() JwtBlackService {
	return newJwtBlackService(s)
}

func (s *service) SysUsers() SysUserService {
	return newUserService(s)
}

// NewService 初始化 Service 接口
func NewService(store store.IStore, rdb *redis.Client) Service {
	return &service{
		store: store,
		rdb:   rdb,
	}
}
