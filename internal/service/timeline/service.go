package timeline

import (
	"context"
	"time"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

//go:generate mockgen -source=service.go -destination=mocks/timeline_mocks.go -package=mocks

type StorageRepo interface {
	SelectFollowersByUserID(ctx context.Context, userID string) ([]string, error)
	SelectTweetsByTweetsIDs(ctx context.Context, tweetIDs []string) ([]domain.Tweet, error)
	SelectLastTweetsByUsersID(ctx context.Context, userIDs []string) ([]domain.Tweet, error)
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	LPush(ctx context.Context, key string, values ...interface{}) error
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
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
