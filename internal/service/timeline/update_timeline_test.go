package timeline_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/renzonaitor/tweet-api/internal/service/timeline"
	"github.com/renzonaitor/tweet-api/internal/service/timeline/mocks"
	"go.uber.org/mock/gomock"
)

func TestUpdateTimeline(t *testing.T) {
	// Define reusable variables for tests
	authorID := uuid.NewString()
	tweetID := uuid.NewString()
	follower1 := uuid.NewString()
	follower2 := uuid.NewString()
	followers := []string{follower1, follower2}

	// Define reusable errors
	dbError := errors.New("database connection failed")
	cacheError := errors.New("redis command failed")

	testCases := []struct {
		name       string
		authorID   string
		tweetID    string
		setupMocks func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository)
	}{
		{
			name:     "Success - Fan-out to multiple followers",
			authorID: authorID,
			tweetID:  tweetID,
			setupMocks: func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository) {
				// Expect storage to be called to get followers, and it succeeds.
				storage.EXPECT().
					SelectFollowersByUserID(gomock.Any(), authorID).
					Return(followers, nil).
					Times(1)

				// Expect cache to be called for each follower.
				timelineKey1 := fmt.Sprintf("timeline:%s", follower1)
				cache.EXPECT().
					LPush(gomock.Any(), timelineKey1, tweetID).
					Return(nil).
					Times(1)

				timelineKey2 := fmt.Sprintf("timeline:%s", follower2)
				cache.EXPECT().
					LPush(gomock.Any(), timelineKey2, tweetID).
					Return(nil).
					Times(1)
			},
		},
		{
			name:     "Success - User has no followers",
			authorID: authorID,
			tweetID:  tweetID,
			setupMocks: func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository) {
				// Expect storage to be called, returning an empty slice.
				storage.EXPECT().
					SelectFollowersByUserID(gomock.Any(), authorID).
					Return([]string{}, nil).
					Times(1)

				// Cache should NOT be called if there are no followers.
				cache.EXPECT().LPush(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:     "Failure - Storage error when fetching followers",
			authorID: authorID,
			tweetID:  tweetID,
			setupMocks: func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository) {
				// Expect storage to be called, and it returns an error.
				storage.EXPECT().
					SelectFollowersByUserID(gomock.Any(), authorID).
					Return(nil, dbError).
					Times(1)

				// Cache should NOT be called if storage fails.
				cache.EXPECT().LPush(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:     "Partial Failure - Cache error for one follower",
			authorID: authorID,
			tweetID:  tweetID,
			setupMocks: func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository) {
				// Expect storage to succeed.
				storage.EXPECT().
					SelectFollowersByUserID(gomock.Any(), authorID).
					Return(followers, nil).
					Times(1)

				// Expect LPush for the first follower to FAIL.
				timelineKey1 := fmt.Sprintf("timeline:%s", follower1)
				cache.EXPECT().
					LPush(gomock.Any(), timelineKey1, tweetID).
					Return(cacheError).
					Times(1)

				// IMPORTANT: Expect LPush for the second follower to still be called.
				// This verifies the loop continues on error.
				timelineKey2 := fmt.Sprintf("timeline:%s", follower2)
				cache.EXPECT().
					LPush(gomock.Any(), timelineKey2, tweetID).
					Return(nil).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockStorage := mocks.NewMockStorageRepo(ctrl)
			mockCache := mocks.NewMockCacheRepository(ctrl)

			if tc.setupMocks != nil {
				tc.setupMocks(mockStorage, mockCache)
			}

			service := timeline.NewService(mockStorage, mockCache)

			// Act
			// Since the method is designed to be async and logs errors instead of returning them,
			// we call it directly. The test's assertions are handled by gomock's expectations.
			service.UpdateTimeline(context.Background(), tc.authorID, tc.tweetID)
		})
	}
}
