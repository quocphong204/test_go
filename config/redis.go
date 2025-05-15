package config

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Port mặc định Redis
		Password: "",               // Không có password
		DB:       0,                // Dùng DB số 0
	})
}
