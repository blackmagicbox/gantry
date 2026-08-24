package server

import (
	"context"
	"sync"

	duelv1 "github.com/blackmagicbox/gantry/gen/go/gantry/duel/v1"
)

type MatchesMap struct {
	mu      sync.Mutex
	matches map[string]string
}

type DuelServer struct {
	duelv1.UnimplementedDuelServiceServer
}

func (ds *DuelServer) TriggerDuel(ctx context.Context, req *duelv1.TriggerDuelRequest) (*duelv1.TriggerDuelResponse, error) {
	// TODO: Implement logic later.
	return &duelv1.TriggerDuelResponse{MatchId: "fake_id-1234"}, nil
}
