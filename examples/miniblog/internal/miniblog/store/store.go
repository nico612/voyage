package store

import (
	"gorm.io/gorm"
	"sync"
)

// IStore 定义了 Store 层需要实现的方法.
type IStore interface {
	Users() UserStore
	Posts() PostStore
	DB() *gorm.DB
}

var (
	// 为了避免实例被重复创建，通常我们需要使用 sync.Once 来确保实例只被初始化一次。
	once sync.Once
	// S 全局变量，方便其它包直接调用已初始化好的 S 实例.
	S *datastore
)

// datastore 是 IStore 的一个具体实现.
type datastore struct {
	db *gorm.DB
}

// 确保 datastore 实现了 IStore 接口.
var _ IStore = (*datastore)(nil)

// NewStore 创建一个 IStore 类型的实例.
func NewStore(db *gorm.DB) *datastore {
	// 确保S只被初始化一次
	once.Do(func() {
		S = &datastore{db: db}
	})
	return S
}

func (ds *datastore) Users() UserStore {
	return newUsers(ds.db)
}

func (ds *datastore) Posts() PostStore {
	return newPosts(ds.db)
}

func (ds *datastore) DB() *gorm.DB {
	return ds.db
}
