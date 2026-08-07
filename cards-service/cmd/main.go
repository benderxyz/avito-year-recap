package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("cards-service started",
			"port", cfg.Port,
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed",
				"error", err,
			)
			os.Exit(1)
		}
	}()

	<-stop
	slog.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	log.Println("cards-service stopped")
}
