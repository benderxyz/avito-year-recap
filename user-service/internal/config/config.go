package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ServerPort          string
	CORSAllowedOrigins  []string
	PostgresHost        string
	PostgresPort        string
	PostgresUser        string
	PostgresPassword    string
	PostgresDatabase    string
	PostgresSSLMode     string
	MigrationsDir       string
}

func Load() Config {
	return Config{
		ServerPort:         getEnv("SERVER_PORT", "8082"),
		CORSAllowedOrigins: getEnvList("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		PostgresHost:       getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:       getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:       getEnv("POSTGRES_USER", "recap"),
		PostgresPassword:   getEnv("POSTGRES_PASSWORD", "recap"),
		PostgresDatabase:   getEnv("POSTGRES_DATABASE", "users"),
		PostgresSSLMode:    getEnv("POSTGRES_SSLMODE", "disable"),
		MigrationsDir:      getEnv("MIGRATIONS_DIR", "migrations"),
	}
}

func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDatabase, c.PostgresSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvList(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	list := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			list = append(list, trimmed)
		}
	}

	if len(list) == 0 {
		return fallback
	}
	return list
}
