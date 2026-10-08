package queue

import (
	"sync"
	"time"

	matchmakingv1 "github.com/blackmagicbox/gantry/gen/go/gantry/matchmaking/v1"
	"github.com/google/uuid"
)

type Queue struct {
	mu    sync.Mutex
	queue map[string]Ticket
}

func New() *Queue {
	return &Queue{
		queue: make(map[string]Ticket),
	}
}

func (q *Queue) Enqueue(playerID string) Ticket {
	ticketID := uuid.New()
	q.mu.Lock()
	defer q.mu.Unlock()

	ticket := Ticket{
		TicketID:  ticketID.String(),
		PlayerID:  playerID,
		MatchID:   "",
		Status:    matchmakingv1.Status_STATUS_WAITING,
		CreatedAt: time.Now(),
	}

	q.queue[ticketID.String()] = ticket
	return ticket
}

func (q *Queue) Get(ticketID string) (Ticket, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	ticket, ok := q.queue[ticketID]
	return ticket, ok
}
