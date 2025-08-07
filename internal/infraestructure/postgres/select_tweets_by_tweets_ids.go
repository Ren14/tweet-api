package postgres

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

// SelectTweetsByTweetsIDs retrieves a slice of Tweets that match the given IDs.
func (r Repository) SelectTweetsByTweetsIDs(ctx context.Context, tweetIDs []string) ([]domain.Tweet, error) {
	if len(tweetIDs) == 0 {
		return []domain.Tweet{}, nil
	}

	query := `
		SELECT id, user_id, content, created_at
		FROM tweets
		WHERE id = ANY($1)
		ORDER BY created_at DESC
	`

	// QueryContext is used because we expect multiple rows in the result.
	rows, err := r.db.QueryContext(ctx, query, tweetIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tweets := make([]domain.Tweet, 0, len(tweetIDs))

	for rows.Next() {
		var tweet domain.Tweet
		if err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Text, &tweet.CreatedAt); err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tweets, nil
}
