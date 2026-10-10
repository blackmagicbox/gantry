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
}

func New() *Queue {
	return &Queue{
		sync.Mutex{},
		make(map[string]Ticket),
		make(map[string]string),
	}
}

func (q *Queue) CheckByPlayer(playerID string) bool {
	if _, ok := q.byPlayerID[playerID]; ok {
		return true
	}
	return false
}

func (q *Queue) Enqueue(playerID string) Ticket {
	ticketID := uuid.New()
	q.mu.Lock()
	defer q.mu.Unlock()

	if exists := q.CheckByPlayer(playerID); exists {
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
	return ticket
}

func (q *Queue) Get(ticketID string) (Ticket, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	ticket, ok := q.queue[ticketID]
	return ticket, ok
}
