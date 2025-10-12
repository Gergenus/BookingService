package redispkg

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func InitRedisDB(redisURL, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return client
}
