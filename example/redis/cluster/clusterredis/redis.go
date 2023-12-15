package clusterredis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Option struct {
	Addr     []string
	Password string
}

// 结构体方法
func (o *Option) initClusterClient() (client *redis.ClusterClient, err error) {

	client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: o.Addr, //[]string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
	})

	c, cancle := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancle()

	_, err = client.Ping(c).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
