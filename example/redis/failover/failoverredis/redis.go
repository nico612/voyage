package failoverredis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// redis 哨兵连接模式,
//一般来说，哨兵模式用来解决了以下问题：
//
//1. **主服务器故障切换：** 当 Redis 主服务器出现故障时，哨兵会自动检测到并进行切换，将其中一个从服务器升级为新的主服务器。
//2. **监控与通知：** 哨兵负责监控 Redis 主服务器和从服务器的运行状况，一旦发现异常，会及时通知管理员或其他系统，以便进行相应的处理。
//3. **配置管理：** 哨兵允许对 Redis 的配置进行动态调整，包括设置从服务器的优先级、故障转移的超时时间、以及故障切换的条件等。

// Option 定义一个结构体
type Option struct {
	Master   string // 主服务器名
	Addr     []string
	Password string
}

// InitSentinelClient 哨兵模式连接
func (o *Option) InitSentinelClient() (client *redis.Client, err error) {
	client = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    o.Master, // "master"
		SentinelAddrs: o.Addr,   // []string{"x.x.x.x:26379", "xx.xx.xx.xx:26379", "xxx.xxx.xxx.xxx:26379"}
	})

	c, cancle := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancle()

	_, err = client.Ping(c).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
