package mysql

import (
	"fmt"
	mycasbin "github.com/nico612/voyage/internal/adminsrv/casbin"
	"github.com/nico612/voyage/internal/adminsrv/models"
	"github.com/nico612/voyage/internal/adminsrv/store"
	"github.com/nico612/voyage/internal/pkg/options"
	"github.com/nico612/voyage/pkg/db"
	"github.com/nico612/voyage/pkg/errors"
	"github.com/nico612/voyage/pkg/log"
	"github.com/nico612/voyage/pkg/logger"
	"gorm.io/gorm"
	"sync"
)

var (
	mysqlStore store.IStore
	once       sync.Once
)

// GetMySQLStoreOr 创建获取 Mysql store 层实例，应该仅在服务初始化时传入opts，后续只需传入nil获取实例
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

		// 初始化 casbin
		_ = mycasbin.CreateCasbinOr(dbIns)

		// 初始化store层
		mysqlStore = &datasource{db: dbIns}

		err = migrateDatabase(dbIns)
	})

	if mysqlStore == nil || err != nil {
		log.Errorf("failed to get mysql store factory, errors: %s", err.Error())
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, errors: %w", mysqlStore, err)
	}

	return mysqlStore, nil
}

// migrateDatabase 自动表迁移，只会添加缺失的字段，不会删除或改变当前数据。
func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		models.SysUser{},
		models.JwtBlackList{},
		models.SysAuthority{},
		models.SysApi{},
		models.SysBaseMenu{},
		models.SysBaseMenuParameter{},
		models.SysBaseMenuBtn{},
		models.SysAuthorityBtn{},
	)
}

// cleanDatabase 清除数据库表
func cleanDatabase(db *gorm.DB) error {
	err := db.Migrator().DropTable(
		models.SysUser{},
	)

	if err != nil {
		return errors.Wrap(err, "drop sysuser table failed")
	}
	return nil
}

// resetDatabase 重置数据库表
func resetDatabase(db *gorm.DB) error {
	if err := cleanDatabase(db); err != nil {
		return err
	}
	if err := migrateDatabase(db); err != nil {
		return err
	}

	return nil
}
