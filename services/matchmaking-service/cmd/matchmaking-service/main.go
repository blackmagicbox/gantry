package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	health_port, ok := os.LookupEnv("HEALTH_PORT")
	if !ok {
		slog.Error("HEALTH_PORT is not set")
		os.Exit(1)
	} else if health_port == "" {
		health_port = "8080"
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		slog.Error("PORT is not set")
		os.Exit(1)
	} else if port == "" {
		port = "50051"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// Do the Server Stuff
	// Create a Handler (mux)
	mux := http.NewServeMux()
	// handle the '/healtz' endpoint on it
	mux.HandleFunc("/healtz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK\n"))
	})
	// Create a new server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", health_port),
		Handler: mux,
	}
	// Create a go routine to run the server
	go func() {
		slog.Info("Starting matchmaking-service", "port", health_port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start matchmaking-service", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

}
