package interceptors

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type playerIDKey struct{}

func UnaryAuthenticationInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	// Record the start time
	start := time.Now()

	// Extract the Authentication player_id
	playerID, err := playerIDFromMetadata(ctx)
	if err != nil {
		slog.Error("Failed to get Authentication player_id", "error", err)
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Log the incoming request.
	slog.Info("Received unary RPC request", "player_id", playerID, "info", info.FullMethod)
	// Call the Actual handler
	ctx = context.WithValue(ctx, playerIDKey{}, playerID)
	resp, err := handler(ctx, req)

	// Calculate the Duration
	duration := time.Since(start)

	// Extract gRPC status code from the Error
	statusCode := status.Code(err)

	// Log the completion
	slog.Info("Unary RPC request completed", "player_id", playerID, "duration", duration, "status", statusCode)

	// Return response or error
	return resp, err
}

func playerIDFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		slog.Error("It was not possible to get metadata from context", "error", "Metadata is mandatory")
		return "", errors.New("metadata is empty")
	}

	// The x-player-id returns a []string that must be tested before returning an index.
	attr := md.Get("x-player-id")
	if len(attr) == 0 {
		slog.Error("Player ID is missing", "error", "Player ID is mandatory")
		return "", errors.New("player ID is mandatory")
	} else if attr[0] == "" {
		slog.Error("Player ID is empty", "error", "Player ID is mandatory")
		return "", errors.New("player ID is mandatory")

	}

	return attr[0], nil
}

func GetPlayerIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(playerIDKey{}).(string)
	if !ok {
		return "", false
	}
	return id, ok
}
