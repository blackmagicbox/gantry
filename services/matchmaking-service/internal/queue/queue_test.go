package queue

import (
	"fmt"
	"slices"
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
	if len(q.waiting) != 1 {
		t.Errorf("Want 1 waiting player, got %d", len(q.waiting))
	}
}

func TestEnqueueAddsToWaitingInFIFOOrder(t *testing.T) {
	q := New()
	want := []string{"player1", "player2", "player3"}
	for _, p := range want {
		q.Enqueue(p)
	}

	if !slices.Equal(q.waiting, want) {
		t.Errorf("Want waiting %v, got %v", want, q.waiting)
	}
}

func TestEnqueueSamePlayerNotAddedToWaitingTwice(t *testing.T) {
	q := New()
	q.Enqueue("player1")
	q.Enqueue("player1")

	if len(q.waiting) != 1 {
		t.Errorf("Want 1 waiting player, got %d: %v", len(q.waiting), q.waiting)
	}
}

func TestEnqueueConcurrentWaitingCount(t *testing.T) {
	q := New()
	const n = 1000
	var wg sync.WaitGroup
	for i := range n {
		wg.Go(func() {
			q.Enqueue(fmt.Sprintf("player%d", i))
		})
	}
	wg.Wait()

	if len(q.waiting) != n {
		t.Errorf("Want %d waiting players, got %d", n, len(q.waiting))
	}
}

func TestTryMatchWithFewerThanTwoPlayers(t *testing.T) {
	q := New()
	if _, ok := q.TryMatch(); ok {
		t.Errorf("Want false for TryMatch on empty queue, got true")
	}

	q.Enqueue("player1")
	if _, ok := q.TryMatch(); ok {
		t.Errorf("Want false for TryMatch with 1 waiting player, got true")
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
	matchID := match.ID
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

func TestTryMatchWithOnePlayer(t *testing.T) {
	q := New()
	ticket := q.Enqueue("player1")

	if _, ok := q.TryMatch(); ok {
		t.Fatalf("Want false for TryMatch with 1 waiting player, got true")
	}

	// The lone player must be left untouched, still waiting for an opponent
	if !slices.Equal(q.waiting, []string{"player1"}) {
		t.Errorf("Want waiting [player1], got %v", q.waiting)
	}
	updated, ok := q.Get(ticket.TicketID)
	if !ok {
		t.Fatalf("Want true, got false for ticket: %s", ticket.TicketID)
	}
	if updated.Status != matchmakingv1.Status_STATUS_WAITING {
		t.Errorf("Want status WAITING, got %s", updated.Status)
	}
	if updated.MatchID != "" {
		t.Errorf("Want empty MatchID, got %s", updated.MatchID)
	}
}

func TestTryMatchWithThreePlayers(t *testing.T) {
	q := New()
	ticket1 := q.Enqueue("player1")
	ticket2 := q.Enqueue("player2")
	ticket3 := q.Enqueue("player3")

	match, ok := q.TryMatch()
	if !ok {
		t.Fatalf("Want true, got false for TryMatch with 3 waiting players")
	}

	// The two oldest players are matched, in FIFO order
	if len(match.Tickets) != 2 {
		t.Fatalf("Want 2 tickets in match, got %d", len(match.Tickets))
	}
	if match.Tickets[0].TicketID != ticket1.TicketID || match.Tickets[1].TicketID != ticket2.TicketID {
		t.Errorf("Want tickets [%s %s], got [%s %s]",
			ticket1.TicketID, ticket2.TicketID, match.Tickets[0].TicketID, match.Tickets[1].TicketID)
	}

	// The third player stays in the queue, still waiting
	if !slices.Equal(q.waiting, []string{"player3"}) {
		t.Errorf("Want waiting [player3], got %v", q.waiting)
	}
	updated3, ok := q.Get(ticket3.TicketID)
	if !ok {
		t.Fatalf("Want true, got false for ticket: %s", ticket3.TicketID)
	}
	if updated3.Status != matchmakingv1.Status_STATUS_WAITING {
		t.Errorf("Want status WAITING, got %s", updated3.Status)
	}
	if updated3.MatchID != "" {
		t.Errorf("Want empty MatchID, got %s", updated3.MatchID)
	}

	// A second attempt has only one player left, so no match
	if _, ok := q.TryMatch(); ok {
		t.Errorf("Want false for second TryMatch with 1 waiting player, got true")
	}

	// Once a fourth player joins, player3 gets matched with them
	q.Enqueue("player4")
	match2, ok := q.TryMatch()
	if !ok {
		t.Fatalf("Want true, got false for TryMatch after player4 joined")
	}
	if match2.ID == match.ID {
		t.Errorf("Want distinct match IDs, got %s twice", match.ID)
	}
	if match2.Tickets[0].TicketID != ticket3.TicketID {
		t.Errorf("Want player3's ticket %s first, got %s", ticket3.TicketID, match2.Tickets[0].TicketID)
	}
	if len(q.waiting) != 0 {
		t.Errorf("Want empty waiting, got %v", q.waiting)
	}
}
