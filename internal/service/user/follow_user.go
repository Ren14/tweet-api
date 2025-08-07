package user

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

func (s Service) FollowUser(ctx context.Context, followUser domain.FollowUser) error {
	// TODO [technical debt] validate if already exist
	// relation := s.Storage.GetRelation(ctx, followUser) TODO implements this method on repository
	// validate if relation is equal to followUser entity TODO implements this logic on receiver function

	// TODO another approach if capture duplicate_key error and return elegant message to user

	return s.Storage.CreateRelation(ctx, followUser)
}
