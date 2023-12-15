package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func main() {

}

// 常用的字符串操作

// 示例1.redis数据库中字符串的set与get操作实践.
func setGetExample(client *redis.Client, ctx context.Context) {
	// 1. Set 设置 key 过期时间 -1 则表示用不过期
	err := client.Set(ctx, "score", 100, 60*time.Second).Err()
	if err != nil {
		fmt.Printf("set score failed, err: %v\n", err)
		panic(err)
	}

	// 2. Get 获取已存在的 Key 其储存的值
	va1, err := client.Get(ctx, "score").Result()
	if err != nil {
		panic(err)
	}

	fmt.Printf("val1 -> score: %v\n", va1)

	// Get 获取一个不存在的值返回redis.Nil 则说明不存在
	val2, err := client.Get(ctx, "name").Result()
	if err == redis.Nil {
		fmt.Println("[ERROR] - Key [name] not exist")
	} else if err != nil {
		fmt.Printf("get name failed, err:%v\n", err)
		panic(err)
	}

	// Exists() 方法用于检测某个key是否存在
	n, _ := client.Exists(ctx, "name").Result()
	if n > 0 {
		fmt.Println("name key 存在!")
	} else {
		fmt.Println("name key 不存在!")
		client.Set(ctx, "name", "weiyi", 60*time.Second)
	}
	val2, _ = client.Get(ctx, "name").Result()
	fmt.Println("val2 -> name : ", val2)

	// 3.SetNX 仅当键不存在时，设置键的字符串值。并设置其过期时间
	val3, err := client.SetNX(ctx, "username", "weiyigeek", 0).Result()
	if err != nil {
		fmt.Printf("set username failed, err:%v\n", err)
		panic(err)
	}
	fmt.Printf("val3 -> username: %v\n", val3)

}
