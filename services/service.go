package services

import (
	"api-service/ratelimiter"
	"api-service/stores"
	"errors"
)

var ErrRateLimit = errors.New("rate limit exceeded")

type Service struct {
	store stores.StatsStore
	rl    ratelimiter.RateLimiter
}

func NewService(s stores.StatsStore, rl ratelimiter.RateLimiter) *Service {
	return &Service{s, rl}
}

// Handle validates rate limit and increments stats
func (s *Service) HandleRequest(userID string) error {
	if !s.rl.Allow(userID) {
		return ErrRateLimit
	}
	return s.store.Increment(userID)
}

// Stats returns per-user request counts
func (s *Service) Stats() (map[string]int, error) {
	return s.store.GetStats()
}
