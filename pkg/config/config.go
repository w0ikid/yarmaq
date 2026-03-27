package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv   string
	LogLevel string
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Zitadel  ZitadelConfig
	Kafka    KafkaConfig
	Services ServiceConfig
	SMTP     SMTPConfig
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

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type ServiceConfig struct {
	AccountsServiceURL    string
	TransactionServiceURL string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
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

func Load(prefix ...string) Config {
	_ = godotenv.Load()

	var servicePrefix string
	if len(prefix) > 0 {
		servicePrefix = prefix[0]
	}

	return Config{
		AppEnv:   getEnv("APP_ENV", "dev", servicePrefix),
		LogLevel: getEnv("LOG_LEVEL", "dev", servicePrefix),
		HTTP: HTTPConfig{
			Port: getEnv("APP_PORT", "8080", servicePrefix),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost", servicePrefix),
			Port:     getEnv("POSTGRES_PORT", "5432", servicePrefix),
			User:     getEnv("POSTGRES_USER", "postgres", servicePrefix),
			Password: getEnv("POSTGRES_PASSWORD", "postgres", servicePrefix),
			DBName:   getEnv("POSTGRES_DB_NAME", "postgres", servicePrefix),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable", servicePrefix),
		},
		Zitadel: ZitadelConfig{
			Domain:  getEnv("ZITADEL_DOMAIN_BACKEND", "http://zitadel.localhost:8080"),
			API:     getEnv("ZITADEL_API_BACKEND", "zitadel.localhost:8080"),
			KeyPath: getEnv("ZITADEL_KEY_PATH", "path/to/key.json", servicePrefix),
			JWKSURL: getEnv("ZITADEL_JWKS_URL", "http://zitadel.localhost:8080/oauth/v2/keys"),
		},
		Kafka: KafkaConfig{
			Brokers: getEnvSlice("KAFKA_BROKERS", "localhost:9092", servicePrefix),
			Topic:   getEnv("KAFKA_TOPIC", "events", servicePrefix),
		},
		Services: ServiceConfig{
			AccountsServiceURL:    getEnv("ACCOUNTS_SERVICE_URL", "http://localhost:8081"),
			TransactionServiceURL: getEnv("TRANSACTION_SERVICE_URL", "http://localhost:8082"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost", servicePrefix),
			Port:     getEnvInt("SMTP_PORT", "1025", servicePrefix),
			Username: getEnv("SMTP_USERNAME", "", servicePrefix),
			Password: getEnv("SMTP_PASSWORD", "", servicePrefix),
			From:     getEnv("SMTP_FROM", "noreply@yarmaq.local", servicePrefix),
			UseTLS:   getEnvBool("SMTP_USE_TLS", "false", servicePrefix),
		},
	}
}

// SMART getEnv: first check service specific variable (e.g. ACCOUNTS_POSTGRES_HOST),
// if not found, check general variable (e.g. POSTGRES_HOST),
// if not found, return fallback
func getEnv(key, fallback string, prefix ...string) string {
	if len(prefix) > 0 && prefix[0] != "" {
		fullKey := fmt.Sprintf("%s_%s", prefix[0], key)
		if value := os.Getenv(fullKey); value != "" {
			return value
		}
	}

	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvSlice(key, fallback string, prefix ...string) []string {
	value := getEnv(key, fallback, prefix...)
	return strings.Split(value, ",")
}

func getEnvInt(key, fallback string, prefix ...string) int {
	value := getEnv(key, fallback, prefix...)
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return n
}

func getEnvBool(key, fallback string, prefix ...string) bool {
	value := getEnv(key, fallback, prefix...)
	return value == "true" || value == "1"
}
