package postgres

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

// CreateRelation inserts a new follow relationship into the database.
// It records that the 'FollowID' user is now following the 'FollowedID' user.
func (r Repository) CreateRelation(ctx context.Context, followUser domain.FollowUser) error {
	// The SQL query with placeholders ($1, $2) to prevent SQL injection.
	query := `
		INSERT INTO follows (follower_id, following_id)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(ctx, query, followUser.FollowID, followUser.FollowedID)
	if err != nil {
		return err
	}

	return nil
}
