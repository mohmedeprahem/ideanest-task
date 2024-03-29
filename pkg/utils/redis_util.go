package utils

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	appConfig, _ := ReadAppConfig()
	return redis.NewClient(&redis.Options{
		Addr:         appConfig.Redis.Address,
		Password:     appConfig.Redis.Password,
		DB:           0,
	})
}

var RedisClient = NewRedisClient()


func SetRedisValue(key string, value bool, expiration time.Duration) error {
	return RedisClient.Set(context.Background(), key, value, expiration).Err()
}

func GetRedisValue(key string) (string, error) {
	return RedisClient.Get(context.Background(), key).Result()
}

func DeleteRedisValue(key string) error {
	return RedisClient.Del(context.Background(), key).Err()
}
