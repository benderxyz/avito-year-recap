package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "cards-service: OK")
	})

	fmt.Println("cards-service started on :8081")

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
