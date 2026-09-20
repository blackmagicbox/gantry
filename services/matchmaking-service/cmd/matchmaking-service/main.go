package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
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
	// Example of a call to Trigger Duel to test the connection with the server
	// conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	slog.Error("it was not possible to create grpc client", "error", err)
	// 	os.Exit(1)
	// }
	// client := duelv1.NewDuelServiceClient(conn)

	// // Temp Test calling Trigger Duel directly
	// req := &duelv1.TriggerDuelRequest{
	// 	IdempotencyKey: "test",
	// 	Player_1Id:     "player-1",
	// 	Player_2Id:     "player-2",
	// }
	// resp, err := client.TriggerDuel(ctx, req)
	// if err != nil {
	// 	slog.Error("Failed to trigger duel", "error", err)
	// } else {
	// 	slog.Info(
	// 		"duel triggered successfully",
	// 		"match_id", resp.MatchId,
	// 	)
	// }

	<-ctx.Done()
	slog.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped gracefully")
}
