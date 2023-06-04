package redis

import (
	"crudly/config"

	"github.com/go-redis/redis/v8"
)

func NewRedis(config config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: config.RedisHost,
	})
}
