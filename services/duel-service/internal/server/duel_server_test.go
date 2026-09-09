package server

import (
	"context"
	"testing"

	"github.com/google/uuid"

	duelv1 "github.com/blackmagicbox/gantry/gen/go/gantry/duel/v1"
)

func TestTriggerDuelCall(t *testing.T) {
	ds := NewDuelServer()

	req := &duelv1.TriggerDuelRequest{
		IdempotencyKey: "idem-key-1",
		Player_1Id:     "player-1",
		Player_2Id:     "player-2",
	}

	resp, err := ds.TriggerDuel(context.Background(), req)
	if err != nil {
		t.Fatalf("TriggerDuel returned unexpected error: %v", err)
	}

	if _, err := uuid.Parse(resp.MatchId); err != nil {
		t.Errorf("MatchId %q is not a valid UUID: %v", resp.MatchId, err)
	}
}

func TestTriggerDuelRequestWithTheSameIdempotentKey(t *testing.T) {
	ds := NewDuelServer()
	req1 := &duelv1.TriggerDuelRequest{
		IdempotencyKey: "idem-key-1",
		Player_1Id:     "player-1",
		Player_2Id:     "player-2",
	}
	req2 := &duelv1.TriggerDuelRequest{
		IdempotencyKey: "idem-key-1",
		Player_1Id:     "player-1",
		Player_2Id:     "player-2",
	}

	resp, err := ds.TriggerDuel(context.Background(), req1)
	if err != nil {
		t.Fatalf("TriggerDuel returned unexpected error: %v", err)
	}

	if _, err := uuid.Parse(resp.MatchId); err != nil {
		t.Errorf("MatchId %q is not a valid UUID: %v", resp.MatchId, err)
	}
	resp2, err := ds.TriggerDuel(context.Background(), req2)
	if err != nil {
		t.Fatalf("TriggerDuel returned unexpected error: %v", err)
	}

	if _, err := uuid.Parse(resp.MatchId); err != nil {
		t.Errorf("MatchId %q is not a valid UUID: %v", resp.MatchId, err)
	}

	if resp.MatchId != resp2.MatchId {
		t.Errorf("got %q, want %q", resp2.MatchId, resp.MatchId)
	}

}

func TestTriggerDuelRequestWithDifferentIdempotentKeys(t *testing.T) {
	ds := NewDuelServer()
	req1 := &duelv1.TriggerDuelRequest{
		IdempotencyKey: "idem-key-1",
		Player_1Id:     "player-1",
		Player_2Id:     "player-2",
	}
	req2 := &duelv1.TriggerDuelRequest{
		IdempotencyKey: "idem-key-2",
		Player_1Id:     "player-1",
		Player_2Id:     "player-2",
	}

	resp, err := ds.TriggerDuel(context.Background(), req1)
	if err != nil {
		t.Fatalf("TriggerDuel returned unexpected error: %v", err)
	}

	if _, err := uuid.Parse(resp.MatchId); err != nil {
		t.Errorf("MatchId %q is not a valid UUID: %v", resp.MatchId, err)
	}

	resp2, err := ds.TriggerDuel(context.Background(), req2)
	if err != nil {
		t.Fatalf("TriggerDuel returned unexpected error: %v", err)
	}

	if _, err := uuid.Parse(resp2.MatchId); err != nil {
		t.Errorf("MatchId %q is not a valid UUID: %v", resp2.MatchId, err)
	}

	if resp.MatchId == resp2.MatchId {
		t.Errorf("got same MatchId %q for different idempotency keys, want different IDs", resp.MatchId)
	}
}

func TestTriggerDuelEmptyIdempotencyKey(t *testing.T) {
	ds := NewDuelServer()

	req := &duelv1.TriggerDuelRequest{
		IdempotencyKey: "",
		Player_1Id:     "player-1",
		Player_2Id:     "player-2",
	}

	resp, err := ds.TriggerDuel(context.Background(), req)
	if err == nil {
		t.Fatalf("TriggerDuel returned no error for empty idempotency key, got resp: %v", resp)
	}

	if resp != nil {
		t.Errorf("got resp %v, want nil on error", resp)
	}
}
