package redis

import (
	"context"
	"fmt"
)

// In your internal/infraestructure/redis/repository.go

func (r *Repository) LPush(ctx context.Context, key string, values ...interface{}) error {
	// The '...' is used to pass the slice elements as individual arguments
	err := r.Client.LPush(ctx, key, values...).Err()
	if err != nil {
		return fmt.Errorf("failed to LPUSH to key %s in redis: %w", key, err)
	}
	return nil
}
