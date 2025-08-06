package timeline

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

func (s Service) GetTimeline(ctx context.Context, userID string, limit int) ([]domain.Tweet, error) {
	// Get last tweets by limit using LRANGE from Redis database
	// TODO implements get last tweets from redis

	// It "hydrates" these IDs by fetching the full tweet objects from PostgreSQL
	followers, err := s.Storage.SelectFollowersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	tweets, err := s.Storage.SelectLastTweetsByUsersID(ctx, followers)
	if err != nil {
		return nil, err
	}

	// It returns the list of hydrated tweets and the next_cursor for pagination.
	return tweets, nil
}
