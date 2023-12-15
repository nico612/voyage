package main

import (
	"context"
	"fmt"
	"github.com/nico612/voyage/example/redis/single/singleredis"
	"github.com/redis/go-redis/v9"
	"time"
)

func main() {

	// 实例化
	opts := &singleredis.Option{
		Addr:         "127.0.0.1:6379",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
		Password:     "",
		Database:     0,
	}
	ctx, cnacle := context.WithTimeout(context.Background(), 1*time.Second)
	defer cnacle()

	// 初始化连接 single redis 服务端
	client, err := opts.InitSingleRedis(ctx)
	if err != nil {
		panic(err)
	}

	// 测试
	V9Example(client)

	defer client.Close() // 关闭redis连接
}

func V9Example(client *redis.Client) {
	ctx := context.Background()

	// 设置Key, 0 表示不过期
	err := client.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	// 获取存在的Key
	val, err := client.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("value = %s\n", val)

	// 获取不存在的Key
	val2, err := client.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}
