package storage

import (
	"fmt"

	"github.com/go-redis/redis"
)

var (
	RedisManagerIns *redis.Client
)

func initRedis() {
	RedisManagerIns = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", "localhost", 6379),
	})
	pong, err := RedisManagerIns.Ping().Result()
	if err != nil {
		fmt.Printf("redis connect error:%v", err)
	}
	fmt.Printf("redis connect success ponged%", pong)
}
