package config

import (
	"event-registration/internal/common"

	"github.com/spf13/viper"
)

func Load() (*common.Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetDefault("POSTGRES_URL", "postgres://user:password@localhost:5432/event_registration?sslmode=disable")
	viper.SetDefault("REDIS_URL", "localhost:6379")
	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("CACHE_TIMEOUT", "5m")
	viper.SetDefault("REFRESH_JWT_EXPIRATION", 7)
	viper.SetDefault("ACCESS_JWT_EXPIRATION", 15)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg common.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
