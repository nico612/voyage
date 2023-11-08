package singleton

import "sync"

type lazySingleton struct {
}

var lazyIns *lazySingleton
var once sync.Once

func GetLazyIns() *lazySingleton {

	once.Do(func() { // 使用once.Do 确保实例全局只被创建一次，还可确保当同时叉棍见多个动作时，只有一个动作被执行
		lazyIns = &lazySingleton{}
	})
	return lazyIns
}
