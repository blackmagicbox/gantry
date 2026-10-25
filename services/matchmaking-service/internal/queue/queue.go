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
	waiting    []string // playerIDs in FIFO order, oldest first
}

func New() *Queue {
	return &Queue{
		sync.Mutex{},
		make(map[string]Ticket),
		make(map[string]string),
		[]string{},
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
	q.waiting = append(q.waiting, playerID)
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
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.waiting) < 2 {
		return Match{}, false
	}
	// Get the two oldest players in `waiting`
	player1 := q.waiting[0]
	player2 := q.waiting[1]
	q.waiting = q.waiting[2:]
	// find the ticketIDs for each byPlayerID
	ticketIDPlayer1 := q.byPlayerID[player1]
	ticketIDPlayer2 := q.byPlayerID[player2]
	// Generate a MatchID
	matchID := uuid.New()
	ticket := q.queue[ticketIDPlayer1]
	ticket2 := q.queue[ticketIDPlayer2]
	// Get the Tickets, change the status
	ticket.Status = matchmakingv1.Status_STATUS_MATCHED
	ticket2.Status = matchmakingv1.Status_STATUS_MATCHED
	// Update the MatchID in both Tickets
	ticket.MatchID = matchID.String()
	ticket2.MatchID = matchID.String()
	q.queue[ticket.TicketID] = ticket
	q.queue[ticket2.TicketID] = ticket2

	// Return the MatchID
	return Match{
		ID:      matchID,
		Tickets: []Ticket{ticket, ticket2}}, true
}
