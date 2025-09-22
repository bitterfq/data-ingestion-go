package main

import (
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	fmt.Println("Starting server on :8080")
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	fmt.Println("Server listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
