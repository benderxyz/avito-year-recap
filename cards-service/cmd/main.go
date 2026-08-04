package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "cards-service: OK")
	})
	mux.HandleFunc("GET /api/recap/{id}", handler.GetRecap)

	log.Println("cards-service started on :8081")

	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}
