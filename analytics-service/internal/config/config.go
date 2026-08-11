package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort         string
	UserServiceURL     string
	ClickHouseHost     string
	ClickHousePort     int
	ClickHouseUser     string
	ClickHousePassword string
	ClickHouseDatabase string
	MigrationsDir      string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	PostgresSSLMode  string

	PGMigrationsDir  string
	PGSeedsDir       string
	SeedOnStart      bool
	RegistryCacheTTL time.Duration
}

func Load() (Config, error) {
	port, err := strconv.Atoi(getEnv("CLICKHOUSE_PORT", "9000"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid CLICKHOUSE_PORT: %w", err)
	}

	registryTTL, err := strconv.Atoi(getEnv("REGISTRY_CACHE_TTL_SECONDS", "60"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid REGISTRY_CACHE_TTL_SECONDS: %w", err)
	}

	return Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		UserServiceURL:     getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		ClickHouseHost:     getEnv("CLICKHOUSE_HOST", "localhost"),
		ClickHousePort:     port,
		ClickHouseUser:     getEnv("CLICKHOUSE_USER", "default"),
		ClickHousePassword: getEnv("CLICKHOUSE_PASSWORD", ""),
		ClickHouseDatabase: getEnv("CLICKHOUSE_DATABASE", "default"),
		MigrationsDir:      getEnv("MIGRATIONS_DIR", "migrations"),

		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "recap"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "recap"),
		PostgresDatabase: getEnv("POSTGRES_DATABASE", "analytics"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),

		PGMigrationsDir:  getEnv("PG_MIGRATIONS_DIR", "migrations-pg"),
		PGSeedsDir:       getEnv("PG_SEEDS_DIR", "seeds"),
		SeedOnStart:      getEnvBool("SEED_ON_START", false),
		RegistryCacheTTL: time.Duration(registryTTL) * time.Second,
	}, nil
}

func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDatabase,
		c.PostgresSSLMode,
	)
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
