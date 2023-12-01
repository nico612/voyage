package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type Options struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	Engine                string           //数据库引擎，默认InnoDB
	Prefix                string           // 表前缀
	Singular              bool             //是否开启全局禁用复数，true表示开启
	MaxOpenConnections    int              // 最大连接数
	MaxConnectionLifeTime time.Duration    // 空闲连接最大存活时间
	MaxIdleConnections    time.Duration    // 最大空闲连接数
	LogLevel              int              // 日志等级
	Logger                logger.Interface // 日志
}

func New(opts *Options) (*gorm.DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Database,
		true,
		"Local")

	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         255,   // string类型字段默认长度
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	gormConfig := &gorm.Config{
		Logger: opts.Logger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   opts.Prefix,
			SingularTable: opts.Singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig)
	if err != nil {
		return nil, err
	}

	// 设置表引擎值
	db.InstanceSet("gorm:table_options", "ENGINE="+opts.Engine)

	// 设置表的字符集和校对规则
	db.InstanceSet("gorm:table_options", "DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)
	sqlDB.SetConnMaxIdleTime(opts.MaxIdleConnections)

	return db, nil
}
