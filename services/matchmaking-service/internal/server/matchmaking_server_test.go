package server

import (
	"context"
	"testing"

	matchmakingv1 "github.com/blackmagicbox/gantry/gen/go/gantry/matchmaking/v1"
	"github.com/blackmagicbox/gantry/services/matchmaking-service/internal/interceptors"
	"github.com/blackmagicbox/gantry/services/matchmaking-service/internal/queue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestJoinQueueEnqueuePlayers(t *testing.T) {
	q := queue.New()
	srv := NewMatchmakingServer(q)
	ctx := interceptors.WithPlayerID(context.Background(), "player1")
	resp, err := srv.JoinQueue(ctx, &matchmakingv1.JoinQueueRequest{})
	if err != nil {
		t.Errorf("Want nil, got %v", err)
	}
	if resp.TicketId == "" {
		t.Errorf("Want non-zero TicketId, got empty string")
	}
	if tkt, ok := q.Get(resp.TicketId); !ok {
		t.Errorf("Want matchmakingv1.TicketId, got %v", resp.TicketId)
	} else if tkt.PlayerID != "player1" {
		t.Errorf("Want player1, got %v", tkt.PlayerID)
	}

}

func TestJoinQueueWithoutPlayerIDReturnsUnauthenticated(t *testing.T) {

	q := queue.New()
	srv := NewMatchmakingServer(q)
	_, err := srv.JoinQueue(context.Background(), &matchmakingv1.JoinQueueRequest{})
	if err == nil {
		t.Fatalf("Want error, got nil")
	}
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("Want Unauthenticated error, got %v", status.Code(err))
	}

}

func TestGetQueueStatusTicketStatusIsWaiting(t *testing.T) {
	q := queue.New()
	srv := NewMatchmakingServer(q)

	ctx := interceptors.WithPlayerID(context.Background(), "player1")

	// _, err := srv.JoinQueue(ctx, &matchmakingv1.JoinQueueRequest{})
	// if err != nil {
	// 	t.Fatalf("expect err to be nil, but got %v instead", err)
	// }

	resp, err := srv.GetQueueStatus(ctx, &matchmakingv1.GetQueueStatusRequest{})
	if err != nil {
		t.Fatalf("expect err to be nil, but got %v instead", err)
	}

	if resp.Status != matchmakingv1.Status_STATUS_WAITING {
		t.Errorf("expect status to be STATUS_WAITING, but got %v instead", resp.Status)
	}

	if resp.MatchId != "" {
		t.Errorf("expect matchId to be empty, but got %v instead", resp.MatchId)
	}
}
