package ratelimiter

import (
	"sync"
	"time"
)

type MemoryLimiter struct {
	mu       sync.Mutex             // protects shared map
	requests map[string][]time.Time // user -> timestamps
	limit    int
	window   time.Duration
}

func NewMemoryLimiter(limit int, window time.Duration) *MemoryLimiter {
	return &MemoryLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow implements sliding window rate limiting
func (m *MemoryLimiter) Allow(userID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	start := now.Add(-m.window)

	// keep only timestamps within window
	var valid []time.Time
	for _, t := range m.requests[userID] {
		if t.After(start) {
			valid = append(valid, t)
		}
	}

	// reject if limit exceeded
	if len(valid) >= m.limit {
		m.requests[userID] = valid
		return false
	}

	// allow request and store timestamp
	valid = append(valid, now)
	m.requests[userID] = valid
	return true
}
