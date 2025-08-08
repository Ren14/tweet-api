package writer

import (
	"context"
	"testing"

	"github.com/renzonaitor/tweet-api/internal/domain"
	"github.com/stretchr/testify/assert"
)

type UserServiceMock struct {
	FollowUserFunc   func(ctx context.Context, followUser domain.FollowUser) error
	PublishTweetFunc func(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error)
}

func (m *UserServiceMock) FollowUser(ctx context.Context, followUser domain.FollowUser) error {
	return m.FollowUserFunc(ctx, followUser)
}

func (m *UserServiceMock) PublishTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error) {
	return m.PublishTweetFunc(ctx, tweet)

}

func Test_WriteHandler(t *testing.T) {
	type args struct {
		userService UserService
	}

	test := []struct {
		name string
		args args
		want *WriterHandler
	}{
		{
			name: "shoud return a new WriterHandler",
			args: args{
				userService: &UserServiceMock{},
			},
			want: &WriterHandler{
				UserService: &UserServiceMock{},
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.args.userService)
			assert.NotNil(t, handler)
		})
	}
}
