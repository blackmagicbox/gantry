package server

import (
	"context"
	"testing"

	"github.com/google/uuid"

	duelv1 "github.com/blackmagicbox/gantry/gen/go/gantry/duel/v1"
)

func newTriggerDuelRequest(idempotencyKey string) *duelv1.TriggerDuelRequest {
	return &duelv1.TriggerDuelRequest{
		IdempotencyKey: idempotencyKey,
		Player_1Id:     "player-1",
		Player_2Id:     "player-2",
	}
}

func requireValidUUID(t *testing.T, id string) {
	t.Helper()
	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("MatchId %q is not a valid UUID: %v", id, err)
	}
}

func TestTriggerDuelCall(t *testing.T) {
	ds := NewDuelServer()

	resp, err := ds.TriggerDuel(context.Background(), newTriggerDuelRequest("idem-key-1"))
	if err != nil {
		t.Fatalf("TriggerDuel returned unexpected error: %v", err)
	}

	requireValidUUID(t, resp.MatchId)
}

func TestTriggerDuelRequestWithTheSameIdempotencyKey(t *testing.T) {
	ds := NewDuelServer()
	req1 := newTriggerDuelRequest("idem-key-1")
	req2 := newTriggerDuelRequest("idem-key-1")

	resp, err := ds.TriggerDuel(context.Background(), req1)
	if err != nil {
		t.Fatalf("TriggerDuel returned unexpected error: %v", err)
	}
	requireValidUUID(t, resp.MatchId)

	resp2, err := ds.TriggerDuel(context.Background(), req2)
	if err != nil {
		t.Fatalf("TriggerDuel returned unexpected error: %v", err)
	}
	requireValidUUID(t, resp2.MatchId)

	if resp.MatchId != resp2.MatchId {
		t.Errorf("got %q, want %q", resp2.MatchId, resp.MatchId)
	}
}

func TestTriggerDuelRequestWithDifferentIdempotencyKeys(t *testing.T) {
	ds := NewDuelServer()
	req1 := newTriggerDuelRequest("idem-key-1")
	req2 := newTriggerDuelRequest("idem-key-2")

	resp, err := ds.TriggerDuel(context.Background(), req1)
	if err != nil {
		t.Fatalf("TriggerDuel returned unexpected error: %v", err)
	}
	requireValidUUID(t, resp.MatchId)

	resp2, err := ds.TriggerDuel(context.Background(), req2)
	if err != nil {
		t.Fatalf("TriggerDuel returned unexpected error: %v", err)
	}
	requireValidUUID(t, resp2.MatchId)

	if resp.MatchId == resp2.MatchId {
		t.Errorf("got same MatchId %q for different idempotency keys, want different IDs", resp.MatchId)
	}
}

func TestTriggerDuelEmptyIdempotencyKey(t *testing.T) {
	ds := NewDuelServer()

	resp, err := ds.TriggerDuel(context.Background(), newTriggerDuelRequest(""))
	if err == nil {
		t.Fatalf("TriggerDuel returned no error for empty idempotency key, got resp: %v", resp)
	}

	if resp != nil {
		t.Errorf("got resp %v, want nil on error", resp)
	}
}
