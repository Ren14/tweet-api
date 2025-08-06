package user

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

type StorageRepo interface {
	CreateRelation(ctx context.Context, follow domain.FollowUser) error
	CreateTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error)
	SelectTweetByID(ctx context.Context, tweetID string) (*domain.Tweet, error)
}

// Service depends on the interfaces, not concrete types.
type Service struct {
	Storage StorageRepo
}

func NewService(storage StorageRepo) *Service {
	return &Service{
		Storage: storage,
	}
}
