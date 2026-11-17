package queue

import (
	"context"
	"fmt"
	"sync"
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

	ctx := t.Context()

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
func TestMatcherWithMoreThanOneMatch(t *testing.T) {

	players := []string{
		"player1",
		"player2",
		"player3",
		"player4",
	}
	expected_matches := int(len(players) / 2)
	q := New()

	for _, p := range players {
		q.Enqueue(p)
	}

	ctx := t.Context()

	matches := make(chan Match, expected_matches)

	go q.RunMatcher(ctx, 200*time.Millisecond, func(m Match) {
		matches <- m
	})

	// Deadline both matches should arrive in the first tick the latest
	deadline := time.After(350 * time.Millisecond)
	seen := map[string]bool{}

	for i := 0; i < expected_matches; i++ {
		select {
		case m := <-matches:
			if m.ID == uuid.Nil {
				t.Error("expected match ID to be non-nil")
			}
			for _, tk := range m.Tickets {
				if seen[tk.TicketID] {
					t.Errorf("ticket %s was matched more than once", tk.TicketID)
				}
				seen[tk.TicketID] = true
			}

		case <-deadline:
			t.Fatalf("want 2 tickets within the first tick, but got %d instead", i)
		}
	}
	if len(seen) != 4 {
		t.Errorf("expected 4 unique tickets to be matched, got %d", len(seen))
	}
}
func TestMatcherWithConcurrentEnqueue(t *testing.T) {
	// Number of players to enqueue
	players := 100

	// Number of matches (half the number of players)
	number_of_matches := 50

	q := New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	matches := make(chan Match, number_of_matches)

	var wg sync.WaitGroup

	go q.RunMatcher(ctx, 10*time.Millisecond, func(m Match) {
		matches <- m
	})

	for i := range players {
		wg.Go(func() {
			q.Enqueue(fmt.Sprintf("player-%d", i+1))
		})
	}
	wg.Wait()

	deadline := time.After(2 * time.Second)

	seen := map[string]bool{}

	for i := range number_of_matches {
		select {
		case m := <-matches:
			if m.ID == uuid.Nil {
				t.Error("expect uuid not to be Nil")
			}
			for _, tk := range m.Tickets {
				if seen[tk.TicketID] {
					t.Errorf("want ticket: %s to be matched only once, it was matched twice instead", tk.TicketID)
				}
				seen[tk.TicketID] = true
			}
		case <-deadline:
			t.Fatalf("want RunMatcher to match %d players in %d matches, found %d maches instead", players, number_of_matches, i)
		}
	}

	if len(seen) != players {
		t.Fatalf("expect %d players matched, found %d instead", players, len(seen))
	}

}
