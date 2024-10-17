package redis

import (
	"chatgpt-web-service/pkg/config"
	"chatgpt-web-service/pkg/log"
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
)

var REDIS *redis.Client
var mutex sync.Mutex

func InitRedis() {
	if REDIS == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if REDIS == nil {
			rr := redis.NewClient(&redis.Options{
				Addr:         config.GetConfig().Redis.Addr,
				Password:     config.GetConfig().Redis.Password,
				DB:           config.GetConfig().Redis.DB,
				PoolSize:     config.GetConfig().Redis.PoolSize,
				MinIdleConns: config.GetConfig().Redis.MinIdle,
			})
			_, err := REDIS.Ping(context.Background()).Result()
			if err != nil {
				panic(err)
			}
			REDIS = rr
			log.My_log.Infof("%s,IP:%s\n", "初始化成功", config.GetConfig().Redis.Addr)
		}
	}
}
