package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "user-service: OK")
	})

	fmt.Println("user-service started on :8082")

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
