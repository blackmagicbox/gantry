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
	waiting    *List
}

func New() *Queue {
	return &Queue{
		sync.Mutex{},
		make(map[string]Ticket),
		make(map[string]string),
		NewList(),
	}
}

func (q *Queue) Enqueue(playerID string) Ticket {
	ticketID := uuid.New()
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.byPlayerID[playerID]; ok {
		return q.queue[q.byPlayerID[playerID]]
	}

	ticket := Ticket{
		TicketID:  ticketID.String(),
		PlayerID:  playerID,
		MatchID:   "",
		Status:    matchmakingv1.Status_STATUS_WAITING,
		CreatedAt: time.Now(),
	}

	q.queue[ticketID.String()] = ticket
	q.byPlayerID[playerID] = ticket.TicketID
	q.waiting.Enqueue(playerID)
	return ticket
}

func (q *Queue) Get(ticketID string) (Ticket, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	ticket, ok := q.queue[ticketID]
	return ticket, ok
}

// TryMatch This method will be called by a loop until there are more then 1 player waiting in the list
func (q *Queue) TryMatch() (Match, bool) {
	// Get the two oldest players in `waiting`

	// find the ticketIDs for each byPlayerID
	// Generate a MatchID
	// Get the Tickets, change the status
	// Update the MatchID in both Tickets
	// Return the MatchID
	return Match{}, false
}
