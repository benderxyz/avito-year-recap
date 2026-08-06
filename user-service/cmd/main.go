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

	"user-service/internal/api"
	"user-service/internal/config"
	"user-service/internal/db"
	"user-service/internal/users"
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
		cancel()
		log.Fatalf("connect clickhouse: %v", err)
	}
	defer func() {
		_ = client.Close()
	}()

	if err := client.Migrate(ctx, cfg.MigrationsDir); err != nil {
		cancel()
		log.Fatalf("migrate: %v", err)
	}

	repo := users.NewRepository(client.Conn())
	handler := api.NewHandler(repo)

	mux := http.NewServeMux()
	handler.Register(mux)
	mux.HandleFunc("GET /api/profiles", api.GetProfiles)
	mux.HandleFunc("GET /internal/users/{id}", api.GetProfile)

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		fmt.Printf("user-service started on :%s\n", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown server: %v", err)
	}
}
