package server

import (
	"context"

	matchmakingv1 "github.com/blackmagicbox/gantry/gen/go/gantry/matchmaking/v1"
	"github.com/blackmagicbox/gantry/services/matchmaking-service/internal/interceptors"
	"github.com/blackmagicbox/gantry/services/matchmaking-service/internal/queue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MatchmakingServer struct {
	queue *queue.Queue
	matchmakingv1.UnimplementedMatchmakingServiceServer
}

func NewMatchmakingServer(q *queue.Queue) *MatchmakingServer {
	return &MatchmakingServer{
		queue: q,
	}
}

func (ms *MatchmakingServer) JoinQueue(ctx context.Context, req *matchmakingv1.JoinQueueRequest) (*matchmakingv1.JoinQueueResponse, error) {
	playerID, ok := interceptors.PlayerIDFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	ticket := ms.queue.Enqueue(playerID)

	return &matchmakingv1.JoinQueueResponse{TicketId: ticket.TicketID}, nil
}

func (ms *MatchmakingServer) GetQueueStatus(ctx context.Context, req *matchmakingv1.GetQueueStatusRequest) (*matchmakingv1.GetQueueStatusResponse, error) {
	return &matchmakingv1.GetQueueStatusResponse{Status: matchmakingv1.Status_STATUS_MATCHED, MatchId: "matchId1"}, nil
}
