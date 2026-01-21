package common

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	PostgresURL               string
	RedisURL                  string
	RabbitMQURL               string
	CacheTimeout              time.Duration
	ServerAddress             string
	ServerPort                string
	ServerExporterAddress     string
	ServerExporterPort        string
	PostgresPlnMobileURL      string
	PostgresPlnMobileHost     string
	PostgresPlnMobilePort     string
	PostgresPlnMobileDatabase string
	PostgresPlnMobileUser     string
	PostgresPlnMobilePassword string
	SshAddress                string
	SshUsername               string
	SshPassword               string
	IsProduction              bool
	GoogleClientSecret        string
	GoogleClientID            string
	GoogleRedirectUri         string
	GoogleOAuthScope          string
	AuthDB                    string
	AuthDBSchema              string
	AuthDBHost                string
	AuthDBPort                string
	AuthDBUser                string
	AuthDBPassword            string
	GarminDB                  string
	GarminDBSchema            string
	GarminDBHost              string
	GarminDBPort              string
	GarminDBUser              string
	GarminDBPassword          string
	GarminToken               string
	GarminRefreshToken        string
	RedisHost                 string
	RedisPort                 int
	RedisPassword             string
	RedisDB                   int
	JwtSecret                 string
	RefreshTokenExpiration    int
	AccessJwtExpiration       int
	SentryDSN                 string
	VCCDBHost                 string
	VCCDBPort                 string
	VCCDBDatabase             string
	VCCDBUser                 string
	VCCDBPassword             string
	VCCDBSchema               string
	DDService                 string
	DDApiKey                  string
	DDSite                    string
	DDENV                     string
	DDVersion                 string
	DDAentHost                string
	DDTraceAgentPort          string
	Timeout                   int
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func Load() (*Config, error) {
	cfg := &Config{
		PostgresURL:               getEnv("POSTGRES_URL", "postgres://postgres:postgres@postgres:5432/vcc?sslmode=disable"),
		RedisURL:                  getEnv("REDIS_URL", "redis:6379"),
		RabbitMQURL:               getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		CacheTimeout:              getEnvDuration("CACHE_TIMEOUT", 5*time.Minute),
		ServerAddress:             getEnv("SERVER_ADDRESS", ""),
		ServerPort:                getEnv("SERVER_PORT", "5051"),
		ServerExporterAddress:     getEnv("SERVER_EXPORTER_ADDRESS", ""),
		ServerExporterPort:        getEnv("SERVER_EXPORTER_PORT", "9091"),
		PostgresPlnMobileURL:      getEnv("POSTGRES_PLN_MOBILE_URL", ""),
		PostgresPlnMobileHost:     getEnv("POSTGRES_PLN_MOBILE_HOST", ""),
		PostgresPlnMobilePort:     getEnv("POSTGRES_PLN_MOBILE_PORT", ""),
		PostgresPlnMobileDatabase: getEnv("POSTGRES_PLN_MOBILE_DATABASE", ""),
		PostgresPlnMobileUser:     getEnv("POSTGRES_PLN_MOBILE_USER", ""),
		PostgresPlnMobilePassword: getEnv("POSTGRES_PLN_MOBILE_PASSWORD", ""),
		SshAddress:                getEnv("SSH_ADDRESS", ""),
		SshUsername:               getEnv("SSH_USERNAME", ""),
		SshPassword:               getEnv("SSH_PASSWORD", ""),
		IsProduction:              getEnvBool("IS_PRODUCTION", false),
		GoogleClientSecret:        getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleClientID:            getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleRedirectUri:         getEnv("GOOGLE_REDIRECT_URI", ""),
		GoogleOAuthScope:          getEnv("GOOGLE_OAUTH_SCOPE", ""),
		AuthDB:                    getEnv("AUTH_DB", ""),
		AuthDBSchema:              getEnv("AUTH_DB_SCHEMA", ""),
		AuthDBHost:                getEnv("AUTH_DB_HOST", ""),
		AuthDBPort:                getEnv("AUTH_DB_PORT", ""),
		AuthDBUser:                getEnv("AUTH_DB_USER", ""),
		AuthDBPassword:            getEnv("AUTH_DB_PASSWORD", ""),
		GarminDB:                  getEnv("GARMIN_DB", ""),
		GarminDBSchema:            getEnv("GARMIN_DB_SCHEMA", ""),
		GarminDBHost:              getEnv("GARMIN_DB_HOST", "postgres"),
		GarminDBPort:              getEnv("GARMIN_DB_PORT", "5432"),
		GarminDBUser:              getEnv("GARMIN_DB_USER", "postgres"),
		GarminDBPassword:          getEnv("GARMIN_DB_PASSWORD", "postgres"),
		GarminToken:               getEnv("GARMIN_TOKEN", ""),
		GarminRefreshToken:        getEnv("GARMIN_REFRESH_TOKEN", ""),
		RedisHost:                 getEnv("REDIS_HOST", "redis"),
		RedisPort:                 getEnvInt("REDIS_PORT", 6379),
		RedisPassword:             getEnv("REDIS_PASSWORD", ""),
		RedisDB:                   getEnvInt("REDIS_DB", 0),
		JwtSecret:                 getEnv("JWT_SECRET", ""),
		RefreshTokenExpiration:    getEnvInt("REFRESH_JWT_EXPIRATION", 7),
		AccessJwtExpiration:       getEnvInt("ACCESS_JWT_EXPIRATION", 15),
		SentryDSN:                 getEnv("SENTRY_DSN", ""),
		VCCDBHost:                 getEnv("POSTGRES_VCC_HOST", ""),
		VCCDBPort:                 getEnv("POSTGRES_VCC_PORT", ""),
		VCCDBDatabase:             getEnv("POSTGRES_VCC_DATABASE", ""),
		VCCDBUser:                 getEnv("POSTGRES_VCC_USER", ""),
		VCCDBPassword:             getEnv("POSTGRES_VCC_PASSWORD", ""),
		VCCDBSchema:               getEnv("POSTGRES_VCC_SCHEMA", ""),
		DDService:                 getEnv("DD_SERVICE", ""),
		DDApiKey:                  getEnv("DD_API_KEY", ""),
		DDSite:                    getEnv("DD_SITE", ""),
		DDENV:                     getEnv("DD_ENV", ""),
		DDVersion:                 getEnv("DD_VERSION", ""),
		DDAentHost:                getEnv("DD_AGENT_HOST", ""),
		DDTraceAgentPort:          getEnv("DD_TRACE_AGENT_PORT", ""),
		Timeout:                   getEnvInt("TIMEOUT", 30),
	}

	return cfg, nil
}
