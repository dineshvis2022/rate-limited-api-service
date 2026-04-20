package main

import (
	"api-service/config"
	"api-service/controllers"
	"api-service/ratelimiter"
	"api-service/routes"
	"api-service/services"
	"api-service/stores"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] .env not found, using system env")
	}

	// Read required configs
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("[ERROR] PORT is required")
	}

	// Storage method (Default: In-Memory)
	useRedis, _ := strconv.ParseBool(os.Getenv("USE_REDIS"))

	// Initialize application (dependency wiring)
	handler := setupApp(useRedis)

	// Setup routes
	router := routes.SetupRouter(handler)

	log.Println("[INFO] Server starting on port:", port)

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Fatal("[ERROR] Failed to start server:", err)
	}
}

// setupApp wires all dependencies (store, rate limiter, service, handler)
func setupApp(useRedis bool) *controllers.Handler {
	var st stores.StatsStore
	var rl ratelimiter.RateLimiter

	if useRedis {
		// Get redis db Connection
		rdb := config.GetRedisDB()

		// Verify connection
		if _, err := rdb.Ping(context.Background()).Result(); err != nil {
			log.Fatal("[ERROR] Redis connection failed:", err)
		}

		log.Println("[INFO] Using Redis (store + rate limiter)")

		// Redis-based implementations
		st = stores.NewRedisStore(rdb)
		rl = ratelimiter.NewRedisLimiter(rdb, 5, time.Minute)

	} else {
		log.Println("[INFO] Using In-Memory (store + rate limiter)")

		// In-memory implementations
		st = stores.NewMemoryStore()
		rl = ratelimiter.NewMemoryLimiter(5, time.Minute)
	}

	// Business layer
	svc := services.NewService(st, rl)

	// HTTP layer
	return controllers.NewHandler(svc)
}
