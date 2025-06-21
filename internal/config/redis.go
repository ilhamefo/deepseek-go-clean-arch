package config

import (
	"github.com/gofiber/storage/redis"
)

func NewRedisConfig(cfg *Config) *redis.Storage {
	return redis.New(redis.Config{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})
}
