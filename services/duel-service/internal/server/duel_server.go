// Package server implements the gRPC handlers for duel-service, backing
// the DuelService API defined in gantry/duel/v1.
package server

import (
	"context"
	"errors"
	"log/slog"

	duelv1 "github.com/blackmagicbox/gantry/gen/go/gantry/duel/v1"
	"github.com/blackmagicbox/gantry/services/duel-service/internal/container"
)

// matches tracks idempotency keys already seen by TriggerDuel, guarding
// against duplicate match creation on retried requests.

// DuelServer implements the duelv1.DuelServiceServer gRPC interface.
// It embeds UnimplementedDuelServiceServer so new RPCs added to the
// proto in the future don't break the build until implemented here.
type DuelServer struct {
	duelv1.UnimplementedDuelServiceServer
	matches *container.Container
}

func NewDuelServer() *DuelServer {
	return &DuelServer{matches: container.NewContainer()}
}

// TriggerDuel handles a request to start a new duel.
//
// TODO: Implement logic later. It currently always returns a fixed,
// non-unique match ID and does not record the match anywhere.
func (ds *DuelServer) TriggerDuel(ctx context.Context, req *duelv1.TriggerDuelRequest) (*duelv1.TriggerDuelResponse, error) {
	ikey := req.IdempotencyKey
	if ikey == "" {
		slog.Error("Missing idempotencyKey")
		return nil, errors.New("idempotency_key is required")
	}

	// Save the Idempotency key to the Container
	// It should create an UUID key and associate to the idempotent key.
	matchID := ds.matches.Inc(ikey)

	return &duelv1.TriggerDuelResponse{MatchId: matchID.String()}, nil
}
