package writer

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

type UserService interface {
	FollowUser(ctx context.Context, followUser domain.FollowUser) error
	PublishTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error)
}

// WriterHandler depends on the interfaces, not concrete types.
type WriterHandler struct {
	UserService UserService
}

func NewHandler(userService UserService) *WriterHandler {
	return &WriterHandler{
		UserService: userService,
	}
}
