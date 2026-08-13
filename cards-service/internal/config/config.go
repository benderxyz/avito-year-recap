package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                string
	LogLevel            string
	UserServiceURL      string
	AnalyticsServiceURL string

	ShareSigningKey    string
	ShareBaseURL       string
	ProductBaseURL     string
	CORSAllowedOrigins []string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	PostgresSSLMode  string

	MigrationsDir string
	SeedsDir      string
	SeedOnStart   bool

	AdminAPIToken string

	LLMAPIKey    string
	LLMEnabled   bool
	LLMProvider  string
	LLMBaseURL   string
	LLMModel     string
	LLMTimeoutMs int
}

func Load() Config {
	return Config{
		Port:                getEnv("PORT", "8081"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		UserServiceURL:      getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		AnalyticsServiceURL: getEnv("ANALYTICS_SERVICE_URL", "http://localhost:8080"),

		ShareSigningKey: getEnv("SHARE_SIGNING_KEY", "dev-insecure-share-key"),
		ShareBaseURL:    getEnv("SHARE_BASE_URL", "http://localhost:3000"),
		ProductBaseURL:  getEnv("PRODUCT_BASE_URL", "https://www.avito.ru"),
		CORSAllowedOrigins: getEnvList(
			"CORS_ALLOWED_ORIGINS",
			[]string{"http://localhost:3000"},
		),

		PostgresHost:     getEnv("POSTGRES_HOST", ""),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "recap"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "recap"),
		PostgresDatabase: getEnv("POSTGRES_DATABASE", "cards"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),

		MigrationsDir: getEnv("MIGRATIONS_DIR", "migrations"),
		SeedsDir:      getEnv("SEEDS_DIR", "seeds"),
		SeedOnStart:   getEnvBool("SEED_ON_START", false),

		AdminAPIToken: getEnv("ADMIN_API_TOKEN", ""),

		LLMAPIKey:    getEnv("OPENAI_API_KEY", ""),
		LLMEnabled:   getEnvBool("LLM_ENABLED", false),
		LLMProvider:  getEnv("LLM_PROVIDER", "openai"),
		LLMBaseURL:   getEnv("LLM_BASE_URL", "https://openrouter.ai/api/v1"),
		LLMModel:     getEnv("LLM_MODEL", "google/gemma-2-9b-it:free"),
		LLMTimeoutMs: getEnvInt("LLM_TIMEOUT_MS", 30000),
	}
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

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
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

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
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
