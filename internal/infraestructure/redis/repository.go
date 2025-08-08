package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/renzonaitor/tweet-api/cmd/http/config"
)

type Repository struct {
	Client *redis.Client
}

// NewRepository creates and configures a new repository with a Redis connection.
func NewRepository(cfg config.Config) (*Repository, error) {
	// 1. Construct the connection address from your config.
	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)

	// 2. Create a new Redis client with the specified options.
	// `go-redis` manages a connection pool for you.
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB, // Default DB is 0
	})

	// 3. Verify the connection is alive.
	// This is a crucial health check to ensure the application starts correctly.
	// We use a short timeout for the initial connection test.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis.")

	// 4. Return the repository with the active client.
	return &Repository{
		Client: rdb,
	}, nil
}

// Close gracefully closes the Redis client and its connection pool.
func (r *Repository) Close() {
	if err := r.Client.Close(); err != nil {
		log.Printf("Error closing Redis client: %v", err)
	}
}
