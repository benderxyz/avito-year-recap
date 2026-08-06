package config

import "os"

type Config struct {
	Port                string
	UserServiceURL      string
	AnalyticsServiceURL string
}

func Load() Config {
	return Config{
		Port:                getEnv("PORT", "8081"),
		UserServiceURL:      getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		AnalyticsServiceURL: getEnv("ANALYTICS_SERVICE_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
