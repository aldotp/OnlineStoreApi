package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type Redis struct {
	RedisClient *redis.Client
}

func NewRedisClient(viper *viper.Viper) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_HOST") + ":" + viper.GetString("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       0,
	})

	return client
}
