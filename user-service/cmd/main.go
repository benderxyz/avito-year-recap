package main

import (
	"fmt"
	"net/http"
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
	}
}
