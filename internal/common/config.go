package common

import "time"

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
	AuthDB                    string        `mapstructure:"AUTH_DB"`
	AuthDBSchema              string        `mapstructure:"AUTH_DB_SCHEMA"`
	AuthDBHost                string        `mapstructure:"AUTH_DB_HOST"`
	AuthDBPort                string        `mapstructure:"AUTH_DB_PORT"`
	AuthDBUser                string        `mapstructure:"AUTH_DB_USER"`
	AuthDBPassword            string        `mapstructure:"AUTH_DB_PASSWORD"`
	RedisHost                 string        `mapstructure:"REDIS_HOST"`
	RedisPort                 int           `mapstructure:"REDIS_PORT"`
	RedisPassword             string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB                   int           `mapstructure:"REDIS_DB"`
	JwtSecret                 string        `mapstructure:"JWT_SECRET"`
	RefreshTokenExpiration    int           `mapstructure:"REFRESH_JWT_EXPIRATION"`
	AccessJwtExpiration       int           `mapstructure:"ACCESS_JWT_EXPIRATION"`
}
