package store

var store IStore

type IStore interface {
	// Users 后管用户接口
	Users() UserStore
}

// Store 提供包级别的 store
func Store() IStore {
	return store
}

func SetStore(s IStore) {
	store = s
}
