package common

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresURL               string        `mapstructure:"POSTGRES_URL"`
	RedisURL                  string        `mapstructure:"REDIS_URL"`
	RabbitMQURL               string        `mapstructure:"RABBITMQ_URL"`
	CacheTimeout              time.Duration `mapstructure:"CACHE_TIMEOUT"`
	ServerAddress             string        `mapstructure:"SERVER_ADDRESS"`
	ServerPort                string        `mapstructure:"SERVER_PORT"`
	ServerExporterAddress     string        `mapstructure:"SERVER_EXPORTER_ADDRESS"`
	ServerExporterPort        string        `mapstructure:"SERVER_EXPORTER_PORT"`
	PostgresPlnMobileURL      string        `mapstructure:"POSTGRES_PLN_MOBILE_URL"`
	PostgresPlnMobileHost     string        `mapstructure:"POSTGRES_PLN_MOBILE_HOST"`
	PostgresPlnMobilePort     string        `mapstructure:"POSTGRES_PLN_MOBILE_PORT"`
	PostgresPlnMobileDatabase string        `mapstructure:"POSTGRES_PLN_MOBILE_DATABASE"`
	PostgresPlnMobileUser     string        `mapstructure:"POSTGRES_PLN_MOBILE_USER"`
	PostgresPlnMobilePassword string        `mapstructure:"POSTGRES_PLN_MOBILE_PASSWORD"`
	SshAddress                string        `mapstructure:"SSH_ADDRESS"`
	SshUsername               string        `mapstructure:"SSH_USERNAME"`
	SshPassword               string        `mapstructure:"SSH_PASSWORD"`
	IsProduction              bool          `mapstructure:"IS_PRODUCTION"`
	GoogleClientSecret        string        `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleClientID            string        `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleRedirectUri         string        `mapstructure:"GOOGLE_REDIRECT_URI"`
	GoogleOAuthScope          string        `mapstructure:"GOOGLE_OAUTH_SCOPE"`
	AuthDB                    string        `mapstructure:"AUTH_DB"`
	AuthDBSchema              string        `mapstructure:"AUTH_DB_SCHEMA"`
	AuthDBHost                string        `mapstructure:"AUTH_DB_HOST"`
	AuthDBPort                string        `mapstructure:"AUTH_DB_PORT"`
	AuthDBUser                string        `mapstructure:"AUTH_DB_USER"`
	AuthDBPassword            string        `mapstructure:"AUTH_DB_PASSWORD"`
	GarminDB                  string        `mapstructure:"GARMIN_DB"`
	GarminDBSchema            string        `mapstructure:"GARMIN_DB_SCHEMA"`
	GarminDBHost              string        `mapstructure:"GARMIN_DB_HOST"`
	GarminDBPort              string        `mapstructure:"GARMIN_DB_PORT"`
	GarminDBUser              string        `mapstructure:"GARMIN_DB_USER"`
	GarminDBPassword          string        `mapstructure:"GARMIN_DB_PASSWORD"`
	GarminToken               string        `mapstructure:"GARMIN_TOKEN"`
	GarminRefreshToken        string        `mapstructure:"GARMIN_REFRESH_TOKEN"`
	RedisHost                 string        `mapstructure:"REDIS_HOST"`
	RedisPort                 int           `mapstructure:"REDIS_PORT"`
	RedisPassword             string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB                   int           `mapstructure:"REDIS_DB"`
	JwtSecret                 string        `mapstructure:"JWT_SECRET"`
	RefreshTokenExpiration    int           `mapstructure:"REFRESH_JWT_EXPIRATION"`
	AccessJwtExpiration       int           `mapstructure:"ACCESS_JWT_EXPIRATION"`
	SentryDSN                 string        `mapstructure:"SENTRY_DSN"`
	VCCDBHost                 string        `mapstructure:"POSTGRES_VCC_HOST"`
	VCCDBPort                 string        `mapstructure:"POSTGRES_VCC_PORT"`
	VCCDBDatabase             string        `mapstructure:"POSTGRES_VCC_DATABASE"`
	VCCDBUser                 string        `mapstructure:"POSTGRES_VCC_USER"`
	VCCDBPassword             string        `mapstructure:"POSTGRES_VCC_PASSWORD"`
	VCCDBSchema               string        `mapstructure:"POSTGRES_VCC_SCHEMA"`
	DDService                 string        `mapstructure:"DD_SERVICE"`
	DDApiKey                  string        `mapstructure:"DD_API_KEY"`
	DDSite                    string        `mapstructure:"DD_SITE"`
	DDENV                     string        `mapstructure:"DD_ENV"`
	DDVersion                 string        `mapstructure:"DD_VERSION"`
	DDAentHost                string        `mapstructure:"DD_AGENT_HOST"`
	DDTraceAgentPort          string        `mapstructure:"DD_TRACE_AGENT_PORT"`
	Timeout                   int           `mapstructure:"TIMEOUT"` // http client timeout in seconds
}

func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetDefault("POSTGRES_URL", "postgres://user:password@localhost:5432/event_registration?sslmode=disable")
	viper.SetDefault("REDIS_URL", "localhost:6379")
	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("CACHE_TIMEOUT", "5m")
	viper.SetDefault("REFRESH_JWT_EXPIRATION", 7)
	viper.SetDefault("ACCESS_JWT_EXPIRATION", 1)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
