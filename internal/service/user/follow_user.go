package user

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

func (s Service) FollowUser(ctx context.Context, followUser domain.FollowUser) error {
	// TODO validate if already exist
	// relation := s.Storage.GetRelation(ctx, followUser)
	// validate if relation is equal to followUser entity

	return s.Storage.CreateRelation(ctx, followUser)
}
