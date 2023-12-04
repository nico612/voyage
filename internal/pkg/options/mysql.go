package options

import (
	"github.com/nico612/voyage/pkg/db"
	"github.com/spf13/pflag"
	"gorm.io/gorm"
	"time"
)

// MySQLOptions MySQL 配置
type MySQLOptions struct {
	Host                  string        `json:"host,omitempty"                     mapstructure:"host"`
	Username              string        `json:"username,omitempty"                 mapstructure:"username"`
	Password              string        `json:"-"                                  mapstructure:"password"`
	Database              string        `json:"database"                           mapstructure:"database"`
	MaxIdleConnections    time.Duration `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"`
	LogLevel              int           `json:"log-level" mapstructure:"log-level"`
}

func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1:3306",
		Username:              "",
		Password:              "",
		Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              1, // Silent
	}
}

func (o *MySQLOptions) NewClient() (*gorm.DB, error) {
	opts := &db.Options{
		Host:                  o.Host,
		Username:              o.Username,
		Password:              o.Password,
		Database:              o.Database,
		MaxIdleConnections:    o.MaxIdleConnections,
		MaxOpenConnections:    o.MaxOpenConnections,
		MaxConnectionLifeTime: o.MaxConnectionLifeTime,
		LogLevel:              o.LogLevel,
	}

	return db.New(opts)
}

func (o *MySQLOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags 将标志位 添加到指定的 pflag.FlagSet 中
func (o *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "mysql.host", o.Host, ""+
		"MySQL service host address. If left blank, the following related mysql options will be ignored.")

	fs.StringVar(&o.Username, "mysql.username", o.Username, ""+
		"Username for access to mysql service.")

	fs.StringVar(&o.Password, "mysql.password", o.Password, ""+
		"Password for access to mysql, should be used pair with password.")

	fs.StringVar(&o.Database, "mysql.database", o.Database, ""+
		"Database name for the server to use.")

	fs.DurationVar(&o.MaxIdleConnections, "mysql.max-idle-connections", o.MaxIdleConnections, ""+
		"Maximum idle connections allowed to connect to mysql.")

	fs.IntVar(&o.MaxOpenConnections, "mysql.max-open-connections", o.MaxOpenConnections, ""+
		"Maximum open connections allowed to connect to mysql.")

	fs.DurationVar(&o.MaxConnectionLifeTime, "mysql.max-connection-life-time", o.MaxConnectionLifeTime, ""+
		"Maximum connection life time allowed to connect to mysql.")

	fs.IntVar(&o.LogLevel, "mysql.log-mode", o.LogLevel, ""+
		"Specify gorm log level.")
}
