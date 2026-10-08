// Command duel-service runs the duel-service gRPC API alongside an HTTP
// health check endpoint, with graceful shutdown on SIGINT/SIGTERM.
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

	duelv1 "github.com/blackmagicbox/gantry/gen/go/gantry/duel/v1"
	"github.com/blackmagicbox/gantry/services/duel-service/internal/server"
	"google.golang.org/grpc"
)

// main starts duel-service's two listeners — an HTTP health check endpoint
// and the gRPC API — and keeps them running until an interrupt or SIGTERM
// triggers a coordinated graceful shutdown of both.
func main() {
	// HEALTH_PORT/PORT are required so deployment manifests must set them
	// explicitly; an empty value (as opposed to unset) falls back to a
	// sane local default for development.
	healthPort, ok := os.LookupEnv("HEALTH_PORT")
	if !ok {
		slog.Error("HEALTH_PORT is not set.")
		os.Exit(1)
	} else if healthPort == "" {
		healthPort = "8080"
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		slog.Error("PORT is not set.")
		os.Exit(1)
	} else if port == "" {
		port = "50051"
	}

	// ctx is canceled on SIGINT/SIGTERM, which is what unblocks the
	// <-ctx.Done() below and kicks off graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Health check server: a minimal HTTP endpoint for orchestrators
	// (e.g. Kubernetes liveness/readiness probes) that is independent of
	// the gRPC API below.
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte("OK\n"))
		if err != nil {
			slog.Error("Failed to write response", "error", err)
			return
		}
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", healthPort),
		Handler: mux,
	}

	go func() {
		slog.Info("Starting duel-service", "health_port", healthPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start duel-service", "error", err)
			os.Exit(1)
		}
	}()

	// gRPC server: the actual duel-service API, served on its own port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		slog.Error("It was not possible to initialize the service listener")
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	duelv1.RegisterDuelServiceServer(grpcServer, server.NewDuelServer())

	go func() {
		slog.Info("Starting the duel-service gRPC server", "port", port)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Block until a shutdown signal arrives, then drain both servers
	// before exiting so in-flight requests aren't cut off.
	<-ctx.Done()
	slog.Info("Shutdown signal received")

	shutDownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutDownContext); err != nil {
		slog.Error("Graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	grpcServer.GracefulStop()
	slog.Info("Server stopped gracefully")
}
