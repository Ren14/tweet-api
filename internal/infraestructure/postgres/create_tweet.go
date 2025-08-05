package postgres

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

func (r Repository) CreateTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error) {
	// The SQL query with placeholders ($1, $2, etc.) to prevent SQL injection.
	query := `
		INSERT INTO tweets (id, user_id, content, created_at)
		VALUES ($1, $2, $3, $4)
	`

	// `ExecContext` is used for queries that don't return rows (INSERT, UPDATE, DELETE).
	result, err := r.db.ExecContext(ctx, query, tweet.ID, tweet.UserID, tweet.Text, tweet.CreatedAt)
	if err != nil {
		return domain.Tweet{}, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return domain.Tweet{}, err
	}
	if rows != 1 {
		// Add Warning Log
	}

	return tweet, nil
}
