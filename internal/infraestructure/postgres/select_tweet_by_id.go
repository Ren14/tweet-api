package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

func (r Repository) SelectTweetByID(ctx context.Context, tweetID string) (*domain.Tweet, error) {
	query := `
		SELECT id, user_id, content, created_at
		FROM tweets
		WHERE id = $1
	`

	// Use QueryRowContext for fetching a single row. It's more efficient
	// than QueryContext for this use case.
	row := r.db.QueryRowContext(ctx, query, tweetID)

	var tweet domain.Tweet
	err := row.Scan(&tweet.ID, &tweet.UserID, &tweet.Text, &tweet.CreatedAt)
	if err != nil {
		// It's a best practice to check specifically for sql.ErrNoRows.
		// This indicates that the tweet was not found, which is a different
		// scenario from a database connection error or a syntax error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // allow to create new tweet
		}
		// For all other potential errors (e.g., connection issues, malformed data),
		// wrap the error to provide more context.
		return nil, fmt.Errorf("error scanning tweet: %w", err)
	}

	return &tweet, nil
}
