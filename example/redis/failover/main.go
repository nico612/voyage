package main

import "github.com/nico612/voyage/example/redis/failover/failoverredis"

func main() {

	opt := failoverredis.Option{
		Master:   "master",
		Addr:     []string{"127.0.0.1:6379"},
		Password: "",
	}

	client, err := opt.InitSentinelClient()
	if err != nil {
		panic(err)
	}

	defer client.Close()
}
