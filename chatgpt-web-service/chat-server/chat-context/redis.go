package chatcontext

import (
	predis "chatgpt-web-service/pkg/db/redis"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	conn *redis.Client
}

func Newrediscache() ContextCache {
	conn := predis.REDIS
	return &RedisCache{
		conn: conn}
}

func (c *RedisCache) Set(key string, value Message, ttl int) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.conn.SetEX(context.Background(), key, string(v), time.Duration(ttl*int(time.Second))).Result()
	return err
}
func (c *RedisCache) Get(key string) (*Message, error) {
	rst, err := c.conn.Get(context.Background(), key).Result()
	if err != nil {
		return &Message{}, err
	}
	message := &Message{}
	err = json.Unmarshal([]byte(rst), &message)
	return message, err
}
