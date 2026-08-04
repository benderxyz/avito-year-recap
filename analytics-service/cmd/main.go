package main

import (
	"fmt"
	"net/http"

	"analytics-service/internal/api"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "analytics-service: OK")
	})
	mux.HandleFunc("GET /internal/metrics/{id}", api.GetMetrics)

	fmt.Println("analytics-service started on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
