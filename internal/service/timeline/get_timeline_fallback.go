package timeline

import (
	"context"
	"log"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

// getTimelineFallback return []tweets from PostgresSQL
func (s Service) getTimelineFallback(ctx context.Context, userID string) ([]domain.Tweet, error) {
	followers, err := s.Storage.SelectFollowersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	tweets, err := s.Storage.SelectLastTweetsByUsersID(ctx, followers)
	if err != nil {
		return nil, err
	}

	// TODO set found tweets_ids in cache using a go-routine for decoupling principal flow

	// TODO add metric response using fallback pattern.
	log.Printf("INFO: return [%d] tweets from fallback PostgreSQL for user: %s", len(tweets), userID)
	return tweets, nil
}
