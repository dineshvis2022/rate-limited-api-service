package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRedisLimiter(c *redis.Client, limit int, window time.Duration) *RedisLimiter {
	return &RedisLimiter{client: c, limit: limit, window: window}
}

// Lua script ensures all operations are atomic inside Redis
var script = redis.NewScript(`
	local key = KEYS[1]
	local now = tonumber(ARGV[1])
	local window = tonumber(ARGV[2])
	local limit = tonumber(ARGV[3])

	-- remove requests outside time window
	redis.call("ZREMRANGEBYSCORE", key, 0, now - window)

	-- count remaining requests
	local count = redis.call("ZCARD", key)

	-- reject if limit exceeded
	if count >= limit then
		return 0
	end

	-- add current request timestamp
	redis.call("ZADD", key, now, now)

	-- set expiry to avoid memory leaks
	redis.call("EXPIRE", key, math.ceil(window / 1000))

	return 1
	`)

// Allow checks if request is permitted
func (r *RedisLimiter) Allow(userID string) bool {
	res, err := script.Run(
		context.Background(),
		r.client,
		[]string{"rate_limit:" + userID},
		time.Now().UnixMilli(),
		r.window.Milliseconds(),
		r.limit,
	).Int()

	// fail-safe: reject if Redis fails
	return err == nil && res == 1
}
