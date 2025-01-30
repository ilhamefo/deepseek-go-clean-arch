// internal/config/config.go
package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds all the configuration settings for the application
type Config struct {
	PostgresURL  string        `mapstructure:"POSTGRES_URL"`  // PostgreSQL connection URL
	RedisURL     string        `mapstructure:"REDIS_URL"`     // Redis connection URL
	RabbitMQURL  string        `mapstructure:"RABBITMQ_URL"`  // RabbitMQ connection URL
	CacheTimeout time.Duration `mapstructure:"CACHE_TIMEOUT"` // Cache timeout duration
}

// Load loads the configuration using Viper
func Load() (*Config, error) {
	// Initialize Viper
	viper.SetConfigName(".env") // Name of the config file (without extension)
	viper.SetConfigType("env")  // Type of the config file (e.g., env, json, yaml)
	viper.AddConfigPath(".")    // Path to look for the config file

	// Set default values
	viper.SetDefault("POSTGRES_URL", "postgres://user:password@localhost:5432/event_registration?sslmode=disable")
	viper.SetDefault("REDIS_URL", "localhost:6379")
	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("CACHE_TIMEOUT", "5m")

	// Read environment variables
	viper.AutomaticEnv() // Automatically override with environment variables

	// Read the config file (if it exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Ignore if the config file is not found, but return other errors
			return nil, err
		}
	}

	// Unmarshal the configuration into the Config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
