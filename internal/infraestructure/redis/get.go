package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Get retrieves a value from Redis.
func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		// It's good practice to check for redis.Nil to know if the key simply doesn't exist.
		if err == redis.Nil {
			return "", nil // Return empty string and no error if key not found
		}
		return "", fmt.Errorf("failed to get key %s from redis: %w", key, err)
	}
	return val, nil
}
