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
	"cards-service/internal/cards"
	"cards-service/internal/clients"
	"cards-service/internal/config"
	"cards-service/internal/db"
)

func main() {
	cfg := config.Load()

	ruleProvider := loadRuleProvider(cfg)

	handler := api.NewHandler(
		clients.NewUserClient(cfg.UserServiceURL),
		clients.NewAnalyticsClient(cfg.AnalyticsServiceURL),
		cfg.ShareSigningKey,
		cfg.ShareBaseURL,
		cfg.ProductBaseURL,
		ruleProvider,
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

func loadRuleProvider(cfg config.Config) *cards.RuleProvider {
	if cfg.PostgresHost == "" {
		slog.Info("postgres not configured, using built-in rules")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pg, err := db.Connect(ctx, cfg.PostgresDSN())
	if err != nil {
		slog.Error("postgres connect failed", "error", err)
		os.Exit(1)
	}

	if err := pg.Migrate(ctx, cfg.MigrationsDir); err != nil {
		slog.Error("postgres migrate failed", "error", err)
		_ = pg.Close()
		os.Exit(1)
	}

	if cfg.SeedOnStart {
		if err := pg.Seed(ctx, cfg.SeedsDir); err != nil {
			slog.Error("postgres seed failed", "error", err)
			_ = pg.Close()
			os.Exit(1)
		}
		slog.Info("rules seeded from files")
	}

	slog.Info("rules loaded from postgres")
	return cards.NewRuleProvider(cards.NewRuleStore(pg.DB()), time.Minute)
}
