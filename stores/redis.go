package stores

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(c *redis.Client) *RedisStore {
	return &RedisStore{client: c}
}

// Increment uses Redis hash to track per-user counts
func (r *RedisStore) Increment(userID string) error {
	return r.client.HIncrBy(context.Background(), "user_stats", userID, 1).Err()
}

// GetStats fetches all user counts from Redis
func (r *RedisStore) GetStats() (map[string]int, error) {
	data, err := r.client.HGetAll(context.Background(), "user_stats").Result()
	if err != nil {
		return nil, err
	}

	res := make(map[string]int)
	for k, v := range data {
		var val int
		fmt.Sscanf(v, "%d", &val) // convert string → int
		res[k] = val
	}
	return res, nil
}
