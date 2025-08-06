package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/renzonaitor/tweet-api/cmd/http/config"
)

// Repository holds the Redis client. The client manages a pool
// of connections automatically.
type Repository struct {
	Client *redis.Client
}

// NewRepository creates and configures a new repository with a Redis connection.
func NewRepository(cfg config.Config) *Repository {
	// 1. Construct the connection address from your config.
	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)

	// 2. Create a new Redis client with the specified options.
	// `go-redis` manages a connection pool for you.
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password, // No password if empty
		DB:       cfg.Redis.DB,       // Default DB is 0
	})

	// 3. Verify the connection is alive.
	// This is a crucial health check to ensure the application starts correctly.
	// We use a short timeout for the initial connection test.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		// If Redis is not available, the application can't function as expected.
		// It's better to fail fast.
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis.")

	// 4. Return the repository with the active client.
	return &Repository{
		Client: rdb,
	}
}

// Close gracefully closes the Redis client and its connection pool.
func (r *Repository) Close() {
	if err := r.Client.Close(); err != nil {
		log.Printf("Error closing Redis client: %v", err)
	}
}
