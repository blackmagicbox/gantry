package queue

import (
	"time"

	matchmakingv1 "github.com/blackmagicbox/gantry/gen/go/gantry/matchmaking/v1"
)

type Ticket struct {
	TicketID string
	PlayerID string
	// For this instance it was decided to use the Status type from the Matchmaking Generated Code, to avoid duplication
	// and unnecessary complexity
	Status    matchmakingv1.Status
	MatchID   string
	CreatedAt time.Time
}
