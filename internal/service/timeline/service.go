package timeline

import (
	"context"
	"time"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

type StorageRepo interface {
	SelectFollowersByUserID(ctx context.Context, userID string) ([]string, error)
	SelectTweetsByTweetsIDs(ctx context.Context, tweetIDs []string) ([]domain.Tweet, error)
	SelectLastTweetsByUsersID(ctx context.Context, userIDs []string) ([]domain.Tweet, error)
}

// CacheRepository defines the contract for a cache.
// This allows for mocking in tests and decouples services from a specific implementation.
type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	LPush(ctx context.Context, key string, values ...interface{}) error
}

// Service depends on the interfaces, not concrete types.
type Service struct {
	Storage StorageRepo
	Cache   CacheRepository
}

func NewService(storage StorageRepo, cache CacheRepository) *Service {
	return &Service{
		Storage: storage,
		Cache:   cache,
	}
}
