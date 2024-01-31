package redis

import (
	"github.com/YungBenn/tech-shop-microservices/configs"
	rdb "github.com/redis/go-redis/v9"
)

func Connect(env configs.EnvVars) *rdb.Client {
	rdb := rdb.NewClient(&rdb.Options{
		Addr:     env.RedisHost,
		DB:       env.RedisDB,
	})

	return rdb
}