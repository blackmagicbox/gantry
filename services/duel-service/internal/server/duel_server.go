// Package server implements the gRPC handlers for duel-service, backing
// the DuelService API defined in gantry/duel/v1.
package server

import (
	"context"

	duelv1 "github.com/blackmagicbox/gantry/gen/go/gantry/duel/v1"
	"github.com/blackmagicbox/gantry/services/duel-service/internal/container"
)

// matches tracks idempotency keys already seen by TriggerDuel, guarding
// against duplicate match creation on retried requests.
var matches = container.NewContainer()

// DuelServer implements the duelv1.DuelServiceServer gRPC interface.
// It embeds UnimplementedDuelServiceServer so new RPCs added to the
// proto in the future don't break the build until implemented here.
type DuelServer struct {
	duelv1.UnimplementedDuelServiceServer
}

// TriggerDuel handles a request to start a new duel.
//
// TODO: Implement logic later. It currently always returns a fixed,
// non-unique match ID and does not record the match anywhere.
func (ds *DuelServer) TriggerDuel(ctx context.Context, req *duelv1.TriggerDuelRequest) (*duelv1.TriggerDuelResponse, error) {

	return &duelv1.TriggerDuelResponse{MatchId: "fake_id-1234"}, nil
}
