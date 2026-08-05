package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"analytics-service/internal/aggregation"
	"analytics-service/internal/api"
	"analytics-service/internal/config"
	"analytics-service/internal/db"
	"analytics-service/internal/events"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	client, err := db.Connect(
		ctx,
		cfg.ClickHouseHost,
		cfg.ClickHousePort,
		cfg.ClickHouseUser,
		cfg.ClickHousePassword,
		cfg.ClickHouseDatabase,
	)
	if err != nil {
		log.Fatalf("connect clickhouse: %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	if err := client.Migrate(ctx, cfg.MigrationsDir); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	registry := events.NewRegistry()
	ingester := events.NewIngester(client.Conn(), registry)
	querier := aggregation.NewClickHouseQuerier(client.Conn())
	timezones := aggregation.NewClickHouseTimezoneResolver(client.Conn())
	metrics := aggregation.NewService(registry, querier, timezones)
	handler := api.NewHandler(ingester, metrics, timezones)

	mux := http.NewServeMux()
	handler.Register(mux)

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		fmt.Printf("analytics-service started on :%s\n", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = server.Shutdown(shutdownCtx)
}
