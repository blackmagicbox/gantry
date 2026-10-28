package queue

import (
	"sync"
	"time"

	matchmakingv1 "github.com/blackmagicbox/gantry/gen/go/gantry/matchmaking/v1"
	"github.com/google/uuid"
)

type Queue struct {
	mu         sync.Mutex
	queue      map[string]Ticket
	byPlayerID map[string]string
	waiting    []string // ticketIDs in FIFO order, oldest first
}

func New() *Queue {
	return &Queue{
		mu:         sync.Mutex{},
		queue:      make(map[string]Ticket),
		byPlayerID: make(map[string]string),
		waiting:    []string{},
	}
}

func (q *Queue) Enqueue(playerID string) Ticket {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.byPlayerID[playerID]; ok {
		ticket := q.queue[q.byPlayerID[playerID]]
		if ticket.Status == matchmakingv1.Status_STATUS_WAITING {
			return ticket
		}
	}

	ticketID := uuid.New()
	ticket := Ticket{
		TicketID:  ticketID.String(),
		PlayerID:  playerID,
		MatchID:   "",
		Status:    matchmakingv1.Status_STATUS_WAITING,
		CreatedAt: time.Now(),
	}

	q.queue[ticketID.String()] = ticket
	q.byPlayerID[playerID] = ticket.TicketID
	q.waiting = append(q.waiting, ticket.TicketID)
	return ticket
}

// Get This method returns the Ticket for a given ticketID
func (q *Queue) Get(ticketID string) (Ticket, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	ticket, ok := q.queue[ticketID]
	return ticket, ok
}

// TryMatch This method will be called by a loop until there are more than 1 player waiting in the list
func (q *Queue) TryMatch() (Match, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.waiting) < 2 {
		return Match{}, false
	}

	// Get the two oldest tickets in `waiting`
	ticketID1 := q.waiting[0]
	ticketID2 := q.waiting[1]
	q.waiting = q.waiting[2:]

	// Generate a MatchID
	matchID := uuid.New()
	ticket := q.queue[ticketID1]
	ticket2 := q.queue[ticketID2]

	// Get the Tickets, change the status
	ticket.Status = matchmakingv1.Status_STATUS_MATCHED
	ticket2.Status = matchmakingv1.Status_STATUS_MATCHED

	// Update the MatchID in both Tickets
	// We won't remove the tickets from the queue because the match hasn't started yet,
	// and we need to be able to query the status of the ticket.
	ticket.MatchID = matchID.String()
	ticket2.MatchID = matchID.String()

	// Update the queue with the tickets
	q.queue[ticket.TicketID] = ticket
	q.queue[ticket2.TicketID] = ticket2

	// Remove the tickets from byPlayerID
	delete(q.byPlayerID, ticket.PlayerID)
	delete(q.byPlayerID, ticket2.PlayerID)

	// Return the MatchID
	return Match{
		ID: matchID,
		Tickets: []Ticket{
			ticket,
			ticket2,
		},
	}, true
}
