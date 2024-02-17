package redis

import (
	rdb "github.com/redis/go-redis/v9"
)

func Connect(redisHost string, redisDB int) *rdb.Client {
	rdb := rdb.NewClient(&rdb.Options{
		Addr:     redisHost,
		DB:       redisDB,
	})

	return rdb
}