package user

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

func (s Service) PublishTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error) {
	return s.Storage.CreateTweet(ctx, tweet)
}
