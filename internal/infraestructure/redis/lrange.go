package redis

import (
	"context"
	"fmt"
)

// LRange retrieves a range of elements from a list in Redis.
func (r *Repository) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	// LRange returns a slice of strings for the given range.
	result, err := r.Client.LRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to LRANGE from key %s in redis: %w", key, err)
	}
	return result, nil
}
