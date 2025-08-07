package redis

import (
	"context"
	"fmt"
	"time"
)

// TODO for code review, this method will be user when implements validation into UpdateTimeline() from timeline service

// Set stores a value in Redis with an expiration.
func (r *Repository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s in redis: %w", key, err)
	}
	return nil
}
