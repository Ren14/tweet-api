package timeline

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

type StorageRepo interface {
	SelectFollowersByUserID(ctx context.Context, userID string) ([]string, error)
	SelectTweetsByTweetsIDs(ctx context.Context, tweetIDs []string) ([]domain.Tweet, error)
	SelectLastTweetsByUsersID(ctx context.Context, userIDs []string) ([]domain.Tweet, error)
}

type CacheRepo interface {
	// todo define contract
}

type Service struct {
	Storage StorageRepo
	Cache   CacheRepo
}

func NewService(storage StorageRepo, cache CacheRepo) *Service {
	return &Service{
		Storage: storage,
		Cache:   cache,
	}
}
