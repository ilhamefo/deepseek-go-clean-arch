package config

import (
	"event-registration/internal/common"

	"github.com/gofiber/storage/redis"
)

func NewRedisConfig(cfg *common.Config) *redis.Storage {
	return redis.New(redis.Config{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})
}
