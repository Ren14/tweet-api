package reader

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

type TimelineService interface {
	GetTimeline(ctx context.Context, userID string, limit int) ([]domain.Tweet, error)
}

type UserService interface {
}

type ReaderHandler struct {
	Timeline TimelineService
	User     UserService
}

func NewHandler(timeline TimelineService, user UserService) *ReaderHandler {
	return &ReaderHandler{
		Timeline: timeline,
		User:     user,
	}
}
