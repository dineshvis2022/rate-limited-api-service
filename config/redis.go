package config

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func GetRedisDB() *redis.Client {
	// Redis config (fallback to default if not provided)
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return rdb
}
