package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	matchmakingv1 "github.com/blackmagicbox/gantry/gen/go/gantry/matchmaking/v1"
	"github.com/blackmagicbox/gantry/services/matchmaking-service/internal/interceptors"
	"github.com/blackmagicbox/gantry/services/matchmaking-service/internal/server"
	"google.golang.org/grpc"
)

func main() {
	healthPort, ok := os.LookupEnv("HEALTH_PORT")
	if !ok {
		slog.Error("HEALTH_PORT is not set")
		os.Exit(1)
	} else if healthPort == "" {
		healthPort = "8080"
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
		_, err := w.Write([]byte("OK\n"))
		if err != nil {
			slog.Error("Failed to write response", "error", err)
			return
		}
	})
	// Create a new server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", healthPort),
		Handler: mux,
	}
	// Create a go routine to run the server
	go func() {
		slog.Info("Starting matchmaking-service", "port", healthPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start matchmaking-service", "error", err)
			os.Exit(1)
		}
	}()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		slog.Error("It was not possible to initialize the service listener", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.UnaryAuthenticationInterceptor))
	matchmakingv1.RegisterMatchmakingServiceServer(grpcServer, server.NewMatchmakingServer())

	go func() {
		slog.Info("Starting the matchmaking-service gRPC server", "port", port)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("Failed to start matchmaking-service gRPC server", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	grpcServer.GracefulStop()
	slog.Info("Server stopped gracefully")
}
