package ratelimiter

// RateLimiter abstracts rate limiting strategy
type RateLimiter interface {
	Allow(userID string) bool
}
