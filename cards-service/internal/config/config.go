package config

import "os"

type Config struct {
	Port                string
	UserServiceURL      string
	AnalyticsServiceURL string
	ShareSigningKey     string
	ShareBaseURL        string
	ProductBaseURL      string
}

func Load() Config {
	return Config{
		Port:                getEnv("PORT", "8081"),
		UserServiceURL:      getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		AnalyticsServiceURL: getEnv("ANALYTICS_SERVICE_URL", "http://localhost:8080"),
		ShareSigningKey:     getEnv("SHARE_SIGNING_KEY", "dev-insecure-share-key"),
		ShareBaseURL:        getEnv("SHARE_BASE_URL", "http://localhost:3000"),
		ProductBaseURL:      getEnv("PRODUCT_BASE_URL", "https://www.avito.ru"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
