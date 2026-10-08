package interceptors

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func newFromContext(md metadata.MD) context.Context {
	ctx := metadata.NewIncomingContext(context.Background(), md)
	return ctx
}

func TestUnaryAuthenticationInterceptor(t *testing.T) {
	md := metadata.Pairs("x-player-id", "player-1")
	handlerCalled := false
	mockHandler := func(ctx context.Context, req any) (any, error) {
		handlerCalled = true
		return "fake Response", nil
	}
	resp, err := UnaryAuthenticationInterceptor(newFromContext(md), "fake Request", &grpc.UnaryServerInfo{FullMethod: "/fake.Service/Method"}, mockHandler)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if resp != "fake Response" {
		t.Error("Expected `fake response`, but got", resp)
	}
	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
}
func TestUnaryAuthenticationInterceptorWithNoPlayerIDSent(t *testing.T) {
	handlerCalled := false
	mockHandler := func(ctx context.Context, req any) (any, error) {
		handlerCalled = true
		return "fake Response", nil
	}
	_, err := UnaryAuthenticationInterceptor(newFromContext(metadata.MD{}), "fake Request", &grpc.UnaryServerInfo{FullMethod: "/fake.Service/Method"}, mockHandler)

	if handlerCalled {
		t.Error("Expected handler not to be called")
	}

	if err == nil {
		t.Errorf("Expected error if player_id not sent, but got %v", err)
	}
}

func TestUnaryAuthenticationInterceptorWithBlankPlayerIDSent(t *testing.T) {
	md := metadata.Pairs("x-player-id", "")
	handlerCalled := false
	mockHandler := func(ctx context.Context, req any) (any, error) {
		handlerCalled = true
		return "fake Response", nil
	}
	_, err := UnaryAuthenticationInterceptor(newFromContext(md), "fake Request", &grpc.UnaryServerInfo{FullMethod: "/fake.Service/Method"}, mockHandler)

	if handlerCalled {
		t.Error("Expected handler not to be called")
	}

	if err == nil {
		t.Errorf("Expected error if player_id is blank, but got %v", err)
	}
}

func TestInterceptorPassesPlayerIDToHandler(t *testing.T) {
	ctx := newFromContext(metadata.Pairs("x-player-id", "player-1"))

	var gotID string
	var gotOK bool

	handler := func(ctx context.Context, req any) (any, error) {
		gotID, gotOK = PlayerIDFromContext(ctx)
		return nil, nil
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/fake/Method"}

	if _, err := UnaryAuthenticationInterceptor(ctx, nil, info, handler); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if !gotOK || gotID != "player-1" {
		t.Errorf("Got(%q, %v), want (%q, true)", gotID, gotOK, "player-1")
	}
}
