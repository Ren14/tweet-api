package postgres

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

// CreateRelation inserts a new follow relationship into the database.
// It records that the 'FollowID' user is now following the 'FollowedID' user.
func (r Repository) CreateRelation(ctx context.Context, followUser domain.FollowUser) error {
	query := `
		INSERT INTO follows (follower_id, following_id)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(ctx, query, followUser.FollowID, followUser.FollowedID)

	return err // nil or error
}
