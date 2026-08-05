package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
<<<<<<< HEAD
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
=======
	"time"

	"user-service/internal/api"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := fmt.Fprintln(w, "user-service: OK"); err != nil {
			http.Error(w, "failed to write health response", http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("GET /api/profiles", api.GetProfiles)
	mux.HandleFunc("GET /internal/users/{id}", api.GetProfile)

	fmt.Println("user-service started on :8082")

	server := &http.Server{
		Addr:              ":8082",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
>>>>>>> feat/implement-services
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

	repo := users.NewRepository(client.Conn())
	handler := api.NewHandler(repo)

	mux := http.NewServeMux()
	handler.Register(mux)

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
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
	_ = server.Shutdown(shutdownCtx)
}
