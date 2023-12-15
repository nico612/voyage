package v1

import (
	"github.com/nico612/voyage/internal/adminsrv/store"
)

type SysMenuAuthorityService interface {
}

type sysMenuAuthorityService struct {
	store store.IStore
}

func newSysMenuAuthorityService(store store.IStore) SysMenuAuthorityService {
	return &sysMenuAuthorityService{store: store}
}
