package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cards-service/internal/api"
	"cards-service/internal/clients"
)

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	userServiceURL := envOrDefault("USER_SERVICE_URL", "http://localhost:8082")
	analyticsServiceURL := envOrDefault("ANALYTICS_SERVICE_URL", "http://localhost:8080")

	handler := api.NewHandler(
		clients.NewUserClient(userServiceURL),
		clients.NewAnalyticsClient(analyticsServiceURL),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := fmt.Fprintln(w, "cards-service: OK"); err != nil {
			http.Error(w, "failed to write health response", http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("GET /api/recap/{id}", handler.GetRecap)

	log.Println("cards-service started on :8081")

	server := &http.Server{
		Addr:              ":8081",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
