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
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	pg, err := db.Connect(ctx, cfg.PostgresDSN())
	if err != nil {
		cancel()
		log.Fatalf("connect postgres: %v", err)
	}
	defer func() {
		_ = pg.Close()
	}()

	if err := pg.Migrate(ctx, cfg.MigrationsDir); err != nil {
		cancel()
		log.Fatalf("migrate: %v", err)
	}

	repo := users.NewRepository(pg.DB())
	handler := api.NewHandler(repo)

	mux := http.NewServeMux()
	handler.Register(mux)

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           api.WithCORS(mux, cfg.CORSAllowedOrigins),
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
