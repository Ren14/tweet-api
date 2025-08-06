package postgres

import (
	"context"
)

func (r Repository) SelectFollowersByUserID(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT following_id
		FROM follows
		WHERE follower_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var followers []string
	for rows.Next() {
		var followerID string
		if err := rows.Scan(&followerID); err != nil {
			return nil, err
		}

		followers = append(followers, followerID)
	}

	return followers, nil
}
