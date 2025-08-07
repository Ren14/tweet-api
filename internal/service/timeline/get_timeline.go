package timeline

import (
	"context"
	"fmt"
	"log"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

func (s Service) GetTimeline(ctx context.Context, userID string, limit int) ([]domain.Tweet, error) {
	timelineKey := fmt.Sprintf(timelineKeyFormat, userID)

	tweetIDs, err := s.Cache.LRange(ctx, timelineKey, 0, int64(limit-1))
	if err != nil {
		return nil, fmt.Errorf("error fetching timeline from cache: %w", err)
	}

	if len(tweetIDs) > 0 {
		// TODO add metric cache hit. This response round the 4-7 ms on localhost test (using Postman)
		log.Printf("INFO: cache hit for key: %s", timelineKey)

		// "Hydrate" the tweet IDs.
		tweets, err := s.Storage.SelectTweetsByTweetsIDs(ctx, tweetIDs)
		if err != nil {
			return nil, fmt.Errorf("error hydrating tweets from storage: %w", err)
		}

		// TODO add metric response ok using cache-first pattern.
		return tweets, nil
	}

	// TODO add metric cache miss. This response round the 8-10 ms on localhost test (using Postman)
	log.Printf("WARN: cache is empty for key: %s. Getting tweets from fallback PostgreSQL", timelineKey)

	return s.getTimelineFallback(ctx, userID)
}
