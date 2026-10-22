package queue

import (
	"fmt"
	"sync"
	"testing"

	matchmakingv1 "github.com/blackmagicbox/gantry/gen/go/gantry/matchmaking/v1"
)

func TestEnqueue(t *testing.T) {
	q := New()
	ticket := q.Enqueue("player1")
	if ticket.CreatedAt.IsZero() {
		t.Errorf("Want non-zero CreatedAt, got zero for ticket %s", ticket.TicketID)
	}
	if ticket.Status != matchmakingv1.Status_STATUS_WAITING {
		t.Errorf("Want status WAITING, got %s", ticket.Status)
	}
	if ticket.TicketID == "" {
		t.Errorf("Want non-zero TicketID, got %s", ticket.TicketID)
	}

	if ticket.PlayerID != "player1" {
		t.Errorf("Want player1, got %s", ticket.PlayerID)
	}

	if _, ok := q.Get(ticket.TicketID); !ok {
		t.Errorf("Want true, got false for ticket: %s", ticket.TicketID)
	}
}

func TestGetWithInvalidTicketID(t *testing.T) {
	q := New()
	ticket := q.Enqueue("player1")
	retrieved, ok := q.Get(ticket.TicketID)
	if !ok {
		t.Errorf("Want true, got false for valid ticketID: %s", ticket.TicketID)
	} else if ticket != retrieved {
		t.Errorf("Want ticket %s, got %s", ticket.TicketID, retrieved.TicketID)
	}
	if _, ok := q.Get("invalidTicketID"); ok {
		t.Errorf("Want false, got true for ticketID: invalidTicketID")
	}
}

func TestEnqueueGenerateDistinctTicketIDs(t *testing.T) {
	q := New()
	first := q.Enqueue("player1")
	second := q.Enqueue("player2")
	if first.TicketID == second.TicketID {
		t.Errorf("Want distinct TicketIDs, got %s", first.TicketID)
	}
}

func TestQueueConcurrentAccess(t *testing.T) {
	q := New()
	const n = 1000
	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() {
			ticket := q.Enqueue(fmt.Sprintf("player%d", i))

			retrievedTicket, ok := q.Get(ticket.TicketID)
			if !ok {
				t.Errorf("Want true, got false for ticket: %s", ticket.TicketID)
			}
			if retrievedTicket.PlayerID != ticket.PlayerID {
				t.Errorf("Want %s, got %s", ticket.PlayerID, retrievedTicket.PlayerID)
			}

		})
	}
	wg.Wait()
}

func TestEnqueueSamePlayerReturnsSameTicket(t *testing.T) {
	q := New()
	first := q.Enqueue("player1")
	second := q.Enqueue("player1")

	if first.TicketID != second.TicketID {
		t.Errorf("Want same %s for player1, got %s", first.TicketID, second.TicketID)
	}
}

func TestEnqueueSamePlayerReturnsSameTicketConcurrent(t *testing.T) {
	q := New()
	const n = 100
	ids := make([]string, n)
	var wg sync.WaitGroup
	for i := range ids {
		wg.Go(func() {
			ticket := q.Enqueue("player1")
			ids[i] = ticket.TicketID
		})
	}
	wg.Wait()
	for i := range ids {
		if ids[0] != ids[i] {
			t.Errorf("Want %s, got %s", ids[0], ids[i])
			break
		}
	}
}

func TestTryMatch(t *testing.T) {
	q := New()
	ticket1 := q.Enqueue("player1")
	ticket2 := q.Enqueue("player2")

	match, ok := q.TryMatch()
	if !ok {
		t.Fatalf("Want true, got false for TryMatch with 2 waiting players")
	}
	matchID := match.MatchId
	if matchID == "" {
		t.Fatalf("Want non-empty matchID, got empty")
	}
	updated1, ok := q.Get(ticket1.TicketID)
	if !ok {
		t.Fatalf("Want true, got false for ticket: %s", ticket1.TicketID)
	}
	updated2, ok := q.Get(ticket2.TicketID)
	if !ok {
		t.Fatalf("Want true, got false for ticket: %s", ticket2.TicketID)
	}
	if updated1.MatchID != matchID {
		t.Errorf("Want %s, got %s", matchID, updated1.MatchID)
	}
	if updated2.MatchID != matchID {
		t.Errorf("Want %s, got %s", matchID, updated2.MatchID)
	}

	if updated1.Status != matchmakingv1.Status_STATUS_MATCHED {
		t.Errorf("Want status MATCHED, got %s", updated1.Status)
	}
	if updated2.Status != matchmakingv1.Status_STATUS_MATCHED {
		t.Errorf("Want status MATCHED, got %s", updated2.Status)
	}

}
