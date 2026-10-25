package queue

import "github.com/google/uuid"

type Match struct {
	ID      uuid.UUID
	Tickets []Ticket
}
