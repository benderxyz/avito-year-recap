package main

import (
	"log"
	"net/http"
	"time"

	"cards-service/internal/api"
	"cards-service/internal/clients"
	"cards-service/internal/config"
)

func main() {
	cfg := config.Load()

	handler := api.NewHandler(
		clients.NewUserClient(cfg.UserServiceURL),
		clients.NewAnalyticsClient(cfg.AnalyticsServiceURL),
	)

	mux := api.RegisterRoutes(handler)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("cards-service started on :%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
