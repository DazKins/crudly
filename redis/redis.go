package redis

import (
	"crudly/config"
	"crypto/tls"

	"github.com/go-redis/redis/v8"
)

func NewRedis(config config.Config) *redis.Client {
	var tlsConfig *tls.Config

	if config.RedisUseSsl {
		tlsConfig = &tls.Config{}
	}

	return redis.NewClient(&redis.Options{
		Addr:      config.RedisHost + ":" + config.RedisPort,
		Password:  config.RedisPassword,
		Username:  config.RedisUsername,
		TLSConfig: tlsConfig,
	})
}
