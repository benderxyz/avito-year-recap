package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort         string
	ClickHouseHost     string
	ClickHousePort     int
	ClickHouseUser     string
	ClickHousePassword string
	ClickHouseDatabase string
	MigrationsDir      string
}

func Load() (Config, error) {
	port, err := strconv.Atoi(getEnv("CLICKHOUSE_PORT", "9000"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid CLICKHOUSE_PORT: %w", err)
	}

	return Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		ClickHouseHost:     getEnv("CLICKHOUSE_HOST", "localhost"),
		ClickHousePort:     port,
		ClickHouseUser:     getEnv("CLICKHOUSE_USER", "default"),
		ClickHousePassword: getEnv("CLICKHOUSE_PASSWORD", ""),
		ClickHouseDatabase: getEnv("CLICKHOUSE_DATABASE", "default"),
		MigrationsDir:      getEnv("MIGRATIONS_DIR", "migrations"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
