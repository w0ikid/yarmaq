package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv   string
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Zitadel  ZitadelConfig
}

type HTTPConfig struct {
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ZitadelConfig struct {
	Domain  string
	API     string
	KeyPath string
	JWKSURL string
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host,
		p.Port,
		p.User,
		p.Password,
		p.DBName,
		p.SSLMode,
	)
}

func Load() Config {
	_ = godotenv.Load()

	appEnv := getEnv("APP_ENV", "dev")

	return Config{
		AppEnv: appEnv,
		HTTP: HTTPConfig{
			Port: getEnv("APP_PORT", "8080"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnv("POSTGRES_SERVICE_DB", "postgres"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		Zitadel: ZitadelConfig{
			Domain:  getEnv("ZITADEL_DOMAIN_BACKEND", "http://zitadel.localhost:8080"),
			API:     getEnv("ZITADEL_API_BACKEND", "zitadel.localhost:8080"),
			KeyPath: getEnv("ZITADEL_KEY_PATH", "path/to/key.json"),
			JWKSURL: getEnv("ZITADEL_JWKS_URL", "http://zitadel.localhost:8080/oauth/v2/keys"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
