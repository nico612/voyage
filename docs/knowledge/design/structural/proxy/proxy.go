// Package proxy 代理模式
package proxy

import "fmt"

type Seller interface {
	sell(name string)
}

// Station 火车站
type Station struct {
	stock int // 库存
}

func (s *Station) sell(name string) {
	if s.stock > 0 {
		s.stock--
		fmt.Printf("代理点中：%s 买了一张票，剩余：%d \n", name, s.stock)
	} else {
		fmt.Println("票已卖空")
	}
}

// StationProxy 火车站代理点
type StationProxy struct {
	station *Station // 持有一个火车站对象
}

func (proxy *StationProxy) sell(name string) {
	if proxy.station.stock > 0 {
		proxy.station.stock--
		fmt.Printf("代理点中：%s 买了一张票，剩余：%d \n", name, proxy.station.stock)
	} else {
		fmt.Println("票已卖空")
	}
}
