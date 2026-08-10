package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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

const (
	shutdownTimeout   = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("user-service stopped with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	pg, err := db.Connect(ctx, cfg.PostgresDSN())
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer func() {
		_ = pg.Close()
	}()

	if err := pg.Migrate(ctx, cfg.MigrationsDir); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	repo := users.NewRepository(pg.DB())
	handler := api.NewHandler(repo, pg.DB())

	mux := http.NewServeMux()
	handler.Register(mux)

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           api.WithCORS(mux, cfg.CORSAllowedOrigins),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	serverErr := make(chan error, 1)

	go func() {
		slog.Info("user-service started", "port", cfg.ServerPort)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("listen: %w", err)
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	slog.Info("user-service stopped")

	return nil
}

func setupLogger(level string) {
	var parsed slog.Level
	if err := parsed.UnmarshalText([]byte(level)); err != nil {
		parsed = slog.LevelInfo
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parsed})))
}
