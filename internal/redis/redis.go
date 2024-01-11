package redis

import (
	"github.com/YungBenn/tech-shop-microservices/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(env config.EnvVars) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.RedisHost,
		DB:       env.RedisDB,
	})

	return rdb
}