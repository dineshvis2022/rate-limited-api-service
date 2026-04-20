package ratelimiter

import (
	"context"
	"fmt"
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
	local member = ARGV[4]

	-- remove old entries
	redis.call("ZREMRANGEBYSCORE", key, 0, now - window)

	-- count current
	local count = redis.call("ZCARD", key)

	if count >= limit then
		return 0
	end

	-- add unique member
	redis.call("ZADD", key, now, member)

	-- set expiry
	redis.call("EXPIRE", key, math.ceil(window / 1000))

	return 1
	`)

// Allow checks if request is permitted
func (r *RedisLimiter) Allow(userID string) bool {
	now := time.Now().UnixMilli()

	// unique member (timestamp + random)
	member := fmt.Sprintf("%d-%d", now, time.Now().UnixNano())

	res, err := script.Run(
		context.Background(),
		r.client,
		[]string{"rate_limit:" + userID},
		now,
		r.window.Milliseconds(),
		r.limit,
		member,
	).Int()

	return err == nil && res == 1
}
