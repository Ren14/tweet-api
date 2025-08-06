package postgres

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

func (r Repository) SelectLastTweetsByUsersID(ctx context.Context, userIDs []string) ([]domain.Tweet, error) {
	if len(userIDs) == 0 {
		return []domain.Tweet{}, nil
	}

	query := `
		SELECT DISTINCT ON (user_id)
		id, user_id, content, created_at
		FROM tweets
		WHERE user_id = ANY($1)
		ORDER BY user_id, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tweets := make([]domain.Tweet, 0, len(userIDs))

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
