package redis_adapter

import "github.com/go-redis/redis"

func New(c *redis.Options) *RedisAdapter {
	conn := redis.NewClient(c)

	return &RedisAdapter{conn: conn}
}

func NewConn(c *redis.Client) *RedisAdapter {
	return &RedisAdapter{conn: c}
}
