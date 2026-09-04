package container

import (
	"sync"

	"github.com/google/uuid"
)

// Container is a concurrency-safe map of match IDs to their state,
// guarded by mu so it can be shared across concurrent gRPC calls.
type Container struct {
	mu         sync.Mutex
	collection map[string]uuid.UUID
}

// NewContainer returns a Container ready for use.
func NewContainer() *Container {
	return &Container{collection: make(map[string]uuid.UUID)}
}

// Inc returns the UUID associated with k, creating and storing a new one
// the first time k is seen.
func (c *Container) Inc(k string) uuid.UUID {
	c.mu.Lock()
	defer c.mu.Unlock()
	if value, ok := c.collection[k]; ok {
		return value
	}
	// Create and save the UUID.
	id := uuid.New()
	c.collection[k] = id
	return id
}
