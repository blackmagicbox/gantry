package queue

import (
	"context"
	"time"
)

// RunMatcher Run TryMatch as go Routine every given interval
func (q *Queue) RunMatcher(ctx context.Context, interval time.Duration, onMatch func(m Match)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			Match, ok := q.TryMatch()
			if ok {
				onMatch(Match)
			}
		}
	}
}
