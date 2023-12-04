package mysql

import (
	"fmt"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/internal/pkg/options"
	"github.com/nico612/voyage/pkg/db"
	"github.com/nico612/voyage/pkg/logger"
	"gorm.io/gorm"
	"sync"
)

// 提供数据储存功能，实现 store.IStore 接口
type datasource struct {
	db *gorm.DB
}

func (ds *datasource) Users() store.UserStore {
	return newUsers(ds)
}

var (
	mysqlStore store.IStore
	once       sync.Once
)

// GetMySQLStoreOr 获取 Mysql store 层实例
func GetMySQLStoreOr(opts *options.MySQLOptions) (store.IStore, error) {
	if opts == nil && mysqlStore == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			Engine:                "InnoDB",
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			MaxIdleConnections:    opts.MaxIdleConnections,
			LogLevel:              opts.LogLevel,
			Logger:                logger.New(opts.LogLevel),
		}
		dbIns, err = db.New(options)
		mysqlStore = &datasource{db: dbIns}
	})

	if mysqlStore == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, errors: %w", mysqlStore, err)
	}

	return mysqlStore, nil
}
