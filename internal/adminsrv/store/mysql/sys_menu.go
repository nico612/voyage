package mysql

import (
	"github.com/nico612/voyage/internal/adminsrv/store"
	"gorm.io/gorm"
)

type sysMenu struct {
	db *gorm.DB
}

var _ store.SysMenu = (*sysMenu)(nil)
