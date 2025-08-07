package redis

import (
	"context"
	"fmt"
)

// LPush insert elements in the head of a list in Redis.
func (r *Repository) LPush(ctx context.Context, key string, values ...interface{}) error {
	err := r.Client.LPush(ctx, key, values...).Err()
	if err != nil {
		return fmt.Errorf("failed to LPUSH to key %s in redis: %w", key, err)
	}
	return nil
}
