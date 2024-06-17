package conf

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func LoadRedis() {
	time.Sleep(Config.RedisMongoSlowLoading)
	RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", Config.RedisUrl.Hostname(), Config.RedisUrl.Port()),
		DB:   Config.RedisDB,
	})

	if cmd := RedisClient.Set(context.Background(), "-cconn-", "-ok-", time.Second); cmd.Err() != nil {
		Logger.Panicf("RedisConnectionCheck: %+v", cmd.Err())
	}
}
