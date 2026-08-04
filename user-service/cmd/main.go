package main

import (
	"fmt"
	"net/http"

	"user-service/internal/api"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "user-service: OK")
	})
	mux.HandleFunc("GET /api/profiles", api.GetProfiles)
	mux.HandleFunc("GET /internal/users/{id}", api.GetProfile)

	fmt.Println("user-service started on :8082")

	if err := http.ListenAndServe(":8082", mux); err != nil {
		panic(err)
	}
}
