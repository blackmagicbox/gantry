package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		slog.Error("PORT is not set.")
		os.Exit(1)
	} else if port == "" {
		port = "50051"
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK\n"))
	})

	slog.Info("Starting duel-service", "port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
