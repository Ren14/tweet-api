package user

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

func (s Service) PublishTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error) {
	existTweet, err := s.Storage.SelectTweetByID(ctx, tweet.ID)
	if err != nil {
		return domain.Tweet{}, err
	}

	if existTweet != nil {
		return *existTweet, nil
	}

	createTweet, err := s.Storage.CreateTweet(ctx, tweet)
	if err != nil {
		return domain.Tweet{}, err
	}

	// This goroutine is a temporary simulation of an async flow.
	// The final implementation should leverage a message broker like AWS SQS/SNS.
	go s.Timeline.UpdateTimeline(context.Background(), tweet.UserID, tweet.ID)

	return createTweet, nil
}
