package postgres

import (
	"context"
	"fmt"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

func (r Repository) CreateRelation(ctx context.Context, followUser domain.FollowUser) error {
	fmt.Println(fmt.Sprintf("followID: %s, followedID: %s", followUser.FollowID, followUser.FollowedID))
	return nil
}
