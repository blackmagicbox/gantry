package server

import (
	"context"

	matchmakingv1 "github.com/blackmagicbox/gantry/gen/go/gantry/matchmaking/v1"
)

type MatchmakingServer struct {
	matchmakingv1.UnimplementedMatchmakingServiceServer
}

func NewMatchmakingServer() *MatchmakingServer {
	return &MatchmakingServer{}
}

func (ms *MatchmakingServer) JoinQueue(ctx context.Context, req *matchmakingv1.JoinQueueRequest) (*matchmakingv1.JoinQueueResponse, error) {
	return &matchmakingv1.JoinQueueResponse{TicketId: "ticketId1"}, nil
}

func (ms *MatchmakingServer) GetQueueStatus(ctx context.Context, rea *matchmakingv1.GetQueueStatusRequest) (*matchmakingv1.GetQueueStatusResponse, error) {
	return &matchmakingv1.GetQueueStatusResponse{Status: matchmakingv1.Status_STATUS_WAITING}, nil
}
