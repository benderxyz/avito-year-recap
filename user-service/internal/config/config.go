package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort       string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	PostgresSSLMode  string
	MigrationsDir    string
}

func Load() Config {
	return Config{
		ServerPort:       getEnv("SERVER_PORT", "8082"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "recap"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "recap"),
		PostgresDatabase: getEnv("POSTGRES_DATABASE", "users"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		MigrationsDir:    getEnv("MIGRATIONS_DIR", "migrations"),
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
