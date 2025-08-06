package user

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

func (s Service) PublishTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error) {
	// validate idempotency by tweet id
	dbTweet, err := s.Storage.SelectTweetByID(ctx, tweet.ID)
	if err != nil {
		return domain.Tweet{}, err
	}

	if dbTweet != nil {
		return *dbTweet, nil
	}

	return s.Storage.CreateTweet(ctx, tweet)
}
