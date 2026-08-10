package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cards-service/internal/api"
	"cards-service/internal/cards"
	"cards-service/internal/clients"
	"cards-service/internal/config"
	"cards-service/internal/db"
	"cards-service/internal/llm"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg)

	pg := setupPostgres(cfg)
	defer func() {
		if err := pg.Close(); err != nil {
			slog.Error("postgres close failed", "error", err)
		}
	}()

	handler := api.NewHandler(
		clients.NewUserClient(cfg.UserServiceURL),
		clients.NewAnalyticsClient(cfg.AnalyticsServiceURL),
		cfg.ShareSigningKey,
		cfg.ShareBaseURL,
		cfg.ProductBaseURL,
		loadRuleProvider(pg.DB()),
	)

	setupLLM(handler, cfg, pg.DB())

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           api.RegisterRoutes(handler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	go func() {
		slog.Info(
			"cards-service started",
			"port", strings.TrimPrefix(server.Addr, ":"),
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-stop
	slog.Info("shutting down...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}

	slog.Info("cards-service stopped")
}

func setupLogger(cfg config.Config) {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: parseLogLevel(cfg.LogLevel),
			},
		),
	)

	slog.SetDefault(logger)
}

func setupPostgres(cfg config.Config) *db.Postgres {
	if cfg.PostgresHost == "" {
		slog.Error("postgres is required for data dictionary rules")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		60*time.Second,
	)
	defer cancel()

	pg, err := db.Connect(ctx, cfg.PostgresDSN())
	if err != nil {
		slog.Error("postgres connect failed", "error", err)
		os.Exit(1)
	}

	if err := pg.Migrate(ctx, cfg.MigrationsDir); err != nil {
		slog.Error("postgres migrate failed", "error", err)

		if closeErr := pg.Close(); closeErr != nil {
			slog.Error("postgres close failed", "error", closeErr)
		}

		os.Exit(1)
	}

	if cfg.SeedOnStart {
		if err := pg.Seed(ctx, cfg.SeedsDir); err != nil {
			slog.Error("postgres seed failed", "error", err)

			if closeErr := pg.Close(); closeErr != nil {
				slog.Error("postgres close failed", "error", closeErr)
			}

			os.Exit(1)
		}

		slog.Info("rules seeded from files")
	}

	return pg
}

func loadRuleProvider(sqlDB *sql.DB) *cards.RuleProvider {
	slog.Info("rules loaded from postgres")

	return cards.NewRuleProvider(
		cards.NewRuleStore(sqlDB),
		time.Minute,
	)
}

func setupLLM(handler *api.Handler, cfg config.Config, sqlDB *sql.DB) {
	if !cfg.LLMEnabled {
		return
	}

	if cfg.LLMAPIKey == "" {
		slog.Warn(
			"LLM_ENABLED is true but OPENAI_API_KEY is empty; serving base recap",
		)
		return
	}

	svc := loadLLMService(cfg, sqlDB)
	if svc == nil {
		return
	}

	slog.Debug("LLM enrichment service loaded")
	handler.SetLLMService(svc)
}

func loadLLMService(
	cfg config.Config,
	sqlDB *sql.DB,
) *llm.Service {
	slog.Debug(
		"loading llm enrichment service",
		"provider", cfg.LLMProvider,
		"model", cfg.LLMModel,
		"timeout_ms", cfg.LLMTimeoutMs,
	)

	timeout := time.Duration(cfg.LLMTimeoutMs) * time.Millisecond

	llmClient := clients.NewLLMClient(
		cfg.LLMAPIKey,
		cfg.LLMBaseURL,
		cfg.LLMModel,
		&http.Client{
			Timeout: timeout + time.Second,
		},
	)

	var provider llm.Provider

	switch cfg.LLMProvider {
	case "", "openai":
		provider = llm.NewOpenAIProvider(llmClient)

	default:
		slog.Error(
			"unknown llm provider, disabling llm",
			"provider", cfg.LLMProvider,
		)
		return nil
	}

	slog.Info(
		"llm enrichment enabled",
		"provider", provider.Name(),
		"model", cfg.LLMModel,
	)

	return llm.NewService(
		provider,
		llm.NewCache(sqlDB),
		timeout,
		cfg.LLMModel,
	)
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
