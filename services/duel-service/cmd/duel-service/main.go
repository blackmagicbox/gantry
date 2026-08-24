package main

import (
	"context"
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

func main() {
	health_port, ok := os.LookupEnv("HEALTH_PORT")
	if !ok {
		slog.Error("HEALTH_PORT is not set.")
		os.Exit(1)
	} else if health_port == "" {
		health_port = "8080"
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		slog.Error("PORT is not set.")
		os.Exit(1)
	} else if port == "" {
		port = "50051"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK\n"))
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", health_port),
		Handler: mux,
	}

	go func() {
		slog.Info("Starting duel-service", "health_port", health_port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		slog.Error("It was not possible to initialize the service listener")
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	duelv1.RegisterDuelServiceServer(grpcServer, &server.DuelServer{})

	go func() {
		// run the newly implemented TCP server
		slog.Info("Starting the duel-service gRPC server", "port", port)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutdown signal received")

	shutDownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutDownContext); err != nil {
		slog.Error("Gracefull shutdown failed", "error", err)
		os.Exit(1)
	}

	grpcServer.GracefulStop()
	slog.Info("Server stopped gracefully")
}
