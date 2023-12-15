package singleredis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// 简单连接 redis 示例

// Option 定义一个Option结构体
type Option struct {
	Addr         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	Password     string
	Database     int
}

// InitSingleRedis 结构体InitSingleRedis方法: 用于初始化redis数据库
func (o *Option) InitSingleRedis(ctx context.Context) (client *redis.Client, err error) {
	// Option

	// Redis 连接对象: NewClient将客户端返回到由选项指定的Redis服务器。
	client = redis.NewClient(&redis.Options{
		Addr:         o.Addr,     // redis服务ip:port
		Password:     o.Password, // redis的认证密码
		WriteTimeout: o.WriteTimeout,
		ReadTimeout:  o.ReadTimeout,
		DialTimeout:  2 * time.Second, // 连接超时时间
		PoolSize:     10,              // 连接池
	})
	fmt.Printf("Connecting Redis : %v\n", o.Addr)
	// 验证是否连接到redis服务端
	res, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Connect Failed! Err: %v\n", err)
		return nil, err
	}

	fmt.Printf("Connect Successful! Ping => %v\n", res)
	return client, nil

}
