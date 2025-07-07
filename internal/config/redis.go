package config

import (
	"event-registration/internal/common"

	"github.com/gofiber/storage/redis"
	redisClient "github.com/redis/go-redis/v9"
)

func NewRedisCache(cfg *common.Config) *redisClient.Client {
	return redisClient.NewClient(&redisClient.Options{
		Addr: cfg.RedisURL,
	})
}

func NewRedisConfig(cfg *common.Config) *redis.Storage {
	return redis.New(redis.Config{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})
}
