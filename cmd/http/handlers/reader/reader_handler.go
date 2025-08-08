package reader

import (
	"context"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

//go:generate mockgen -source=reader_handler.go -destination=./../mocks/timeline_service_mock.go -package=mocks
type TimelineService interface {
	GetTimeline(ctx context.Context, userID string, limit int) ([]domain.Tweet, error)
}

// ReaderHandler depends on the interfaces, not concrete types.
type ReaderHandler struct {
	Timeline TimelineService
}

func NewHandler(timeline TimelineService) *ReaderHandler {
	return &ReaderHandler{
		Timeline: timeline,
	}
}
