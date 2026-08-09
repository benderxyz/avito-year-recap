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
	"cards-service/internal/llm"
)

func main() {

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	slog.SetDefault(logger)
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

	if cfg.LLMEnabled && cfg.LLMAPIKey != "" {
		if svc := loadLLMService(cfg); svc != nil {
			slog.Debug("LLM enrichment service loaded")
			handler.SetLLMService(svc)
		}
	} else if cfg.LLMEnabled {
		slog.Warn("LLM_ENABLED is true but OPENAI_API_KEY is empty; serving base recap")
	}

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
		slog.Error("postgres is required for data dictionary rules")
		os.Exit(1)
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

// loadLLMService builds the optional LLM enrichment service. It uses its own
// Postgres handle for the result cache (the recap_llm_cache table is created by
// the migrations run in loadRuleProvider). Any failure disables LLM by
// returning nil, leaving the base recap flow untouched.
func loadLLMService(cfg config.Config) *llm.Service {
	slog.Debug("loading llm enrichment service",
		"provider", cfg.LLMProvider,
		"model", cfg.LLMModel,
		"timeout_ms", cfg.LLMTimeoutMs,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pg, err := db.Connect(ctx, cfg.PostgresDSN())
	if err != nil {
		slog.Error("llm cache postgres connect failed, disabling llm", "error", err)
		return nil
	}

	timeout := time.Duration(cfg.LLMTimeoutMs) * time.Millisecond

	llmClient := clients.NewLLMClient(
		cfg.LLMAPIKey,
		cfg.LLMBaseURL,
		cfg.LLMModel,
		&http.Client{Timeout: timeout + time.Second},
	)

	var provider llm.Provider
	switch cfg.LLMProvider {
	case "", "openai":
		provider = llm.NewOpenAIProvider(llmClient)
	default:
		slog.Error("unknown llm provider, disabling llm", "provider", cfg.LLMProvider)
		return nil
	}

	slog.Info("llm enrichment enabled", "provider", provider.Name(), "model", cfg.LLMModel)
	return llm.NewService(provider, llm.NewCache(pg.DB()), timeout)
}
