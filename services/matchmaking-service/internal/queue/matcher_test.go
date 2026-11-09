package queue

import (
	"context"
	"testing"
	"time"

	matchmakingv1 "github.com/blackmagicbox/gantry/gen/go/gantry/matchmaking/v1"
	"github.com/google/uuid"
)

func TestMatcher(t *testing.T) {
	q := New()
	_ = q.Enqueue("player1")
	_ = q.Enqueue("Player2")

	matches := make(chan Match, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go q.RunMatcher(ctx, 10*time.Millisecond, func(m Match) {
		matches <- m
	})

	select {
	case m := <-matches:
		if m.ID == uuid.Nil {
			t.Error("expected match ID to be non-nil")
		}
		if len(m.Tickets) != 2 {
			t.Errorf("expected 2 tickets in match, got %d", len(m.Tickets))
		}
		for _, tk := range m.Tickets {
			if tk.Status != matchmakingv1.Status_STATUS_MATCHED {
				t.Errorf("expected ticket status to be MATCHED, got %s", tk.Status)
			}
		}

	case <-time.After(2 * time.Second):
		t.Error("no match received within 2 seconds")
	}
}

func TestMatcherWithCanceledContext(t *testing.T) {
	q := New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})

	go func() {
		q.RunMatcher(ctx, 1000*time.Millisecond, func(m Match) {
		})
		close(done)
	}()
	cancel()

	select {
	case <-done:
		return
	case <-time.After(1 * time.Second):
		t.Fatal("RunMatcher did not return after cancelation")
	}
}
