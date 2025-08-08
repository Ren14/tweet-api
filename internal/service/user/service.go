package user

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

//go:generate mockgen -source=service.go -destination=mocks/user_mocks.go -package=mocks
type StorageRepo interface {
	CreateRelation(ctx context.Context, follow domain.FollowUser) error
	CreateTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error)
	SelectTweetByID(ctx context.Context, tweetID string) (*domain.Tweet, error)
}

type TimelineUpdater interface {
	UpdateTimeline(ctx context.Context, tweetAuthorID, tweetID string)
}

// Service depends on the interfaces, not concrete types.
type Service struct {
	Storage  StorageRepo
	Timeline TimelineUpdater
}

func NewService(storage StorageRepo, timeline TimelineUpdater) *Service {
	return &Service{
		Storage:  storage,
		Timeline: timeline,
	}
}
