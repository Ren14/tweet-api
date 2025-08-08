package reader

import (
	"context"
	"testing"

	"github.com/renzonaitor/tweet-api/internal/domain"
	"github.com/stretchr/testify/assert"
)

type TimelineServiceMock struct {
	GetTimelineFunc func(ctx context.Context, userID string, limit int) ([]domain.Tweet, error)
}

func (m *TimelineServiceMock) GetTimeline(ctx context.Context, userID string, limit int) ([]domain.Tweet, error) {
	return m.GetTimelineFunc(ctx, userID, limit)
}

func Test_NewHandler(t *testing.T) {
	type args struct {
		timelineService TimelineService
	}

	tests := []struct {
		name string
		args args
		want *ReaderHandler
	}{
		{
			name: "should return a new ReaderHandler",
			args: args{
				timelineService: &TimelineServiceMock{},
			},
			want: &ReaderHandler{
				Timeline: &TimelineServiceMock{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				handler := NewHandler(tt.args.timelineService)
				assert.NotNil(t, handler)
			})
		})

	}
}
