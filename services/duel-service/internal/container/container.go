package container

import "sync"

// Container is a concurrency-safe map of match IDs to their state,
// guarded by mu so it can be shared across concurrent gRPC calls.
type Container struct {
	mu         sync.Mutex
	collection map[string]string
}

// NewContainer returns a Container ready for use.
func NewContainer() *Container {
	return &Container{collection: make(map[string]string)}
}

// Inc inserts v under k if k is not already present, returning true if the
// insert happened and false if k was already recorded.
func (c *Container) Inc(k string, v string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.collection[k]; ok {
		return !ok
	}
	c.collection[k] = v
	return true
}
