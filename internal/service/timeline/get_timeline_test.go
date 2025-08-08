package timeline_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/renzonaitor/tweet-api/internal/domain"
	"github.com/renzonaitor/tweet-api/internal/service/timeline"
	"github.com/renzonaitor/tweet-api/internal/service/timeline/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetTimeline(t *testing.T) {
	// Define reusable variables and mock data
	user1 := uuid.NewString()
	user2 := uuid.NewString()
	user3 := uuid.NewString()
	tweet1 := uuid.NewString()
	tweet2 := uuid.NewString()
	tweet3 := uuid.NewString()
	limit := 10
	now := time.Now().Format(time.RFC3339)
	tweetIDs := []string{tweet1, tweet2}
	mockTweets := []domain.Tweet{
		{ID: tweet1, UserID: user1, Text: "Hello from cache!", CreatedAt: now},
		{ID: tweet2, UserID: user2, Text: "This is a test", CreatedAt: now},
	}
	fallbackTweets := []domain.Tweet{
		{ID: tweet2, UserID: user2, Text: "This is a test", CreatedAt: now},
		{ID: tweet3, UserID: user3, Text: "Fallback tweet!", CreatedAt: now},
	}

	fallbackFollowers := []string{user2, user3}

	// Define reusable errors
	cacheError := errors.New("redis connection refused")
	dbError := errors.New("postgres connection failed")

	testCases := []struct {
		name           string
		setupMocks     func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository)
		expectedTweets []domain.Tweet
		expectedErr    error
	}{
		{
			name: "Success - Cache Hit",
			setupMocks: func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository) {
				// 1. Expect a call to the cache, which returns tweet IDs.
				cache.EXPECT().
					LRange(gomock.Any(), gomock.Any(), int64(0), int64(limit-1)).
					Return(tweetIDs, nil).
					Times(1)

				// 2. Expect a call to storage to "hydrate" these IDs.
				storage.EXPECT().
					SelectTweetsByTweetsIDs(gomock.Any(), tweetIDs).
					Return(mockTweets, nil).
					Times(1)
			},
			expectedTweets: mockTweets,
			expectedErr:    nil,
		},
		{
			name: "Success - Cache Miss, Fallback Succeeds",
			setupMocks: func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository) {
				// 1. Expect a call to the cache, which returns an empty slice (cache miss).
				cache.EXPECT().
					LRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]string{}, nil)

				//2. Expect a call to the storage to SelectFollowersByUserID
				storage.EXPECT().SelectFollowersByUserID(gomock.Any(), user1).
					Return(fallbackFollowers, nil)

				// 3. Expect a call to the fallback method in storage.
				storage.EXPECT().
					SelectLastTweetsByUsersID(gomock.Any(), fallbackFollowers).
					Return(fallbackTweets, nil).
					Times(1)
			},
			expectedTweets: fallbackTweets,
			expectedErr:    nil,
		},
		{
			name: "Success - Cache Miss, Fallback is Empty",
			setupMocks: func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository) {
				// 1. Cache miss.
				cache.EXPECT().
					LRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]string{}, nil)

				//2. Expect a call to the storage to SelectFollowersByUserID
				storage.EXPECT().SelectFollowersByUserID(gomock.Any(), user1).
					Return(fallbackFollowers, nil)

				// 3. Fallback returns no tweets.
				storage.EXPECT().
					SelectLastTweetsByUsersID(gomock.Any(), fallbackFollowers).
					Return([]domain.Tweet{}, nil)
			},
			expectedTweets: []domain.Tweet{},
			expectedErr:    nil,
		},
		{
			name: "Failure - Cache Error",
			setupMocks: func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository) {
				// Expect a call to the cache, which returns an error.
				cache.EXPECT().
					LRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, cacheError)
				// No calls to storage should be made if the cache fails.
			},
			expectedTweets: nil,
			expectedErr:    cacheError,
		},
		{
			name: "Failure - Cache Hit, Hydration Fails",
			setupMocks: func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository) {
				// 1. Cache returns IDs successfully.
				cache.EXPECT().
					LRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tweetIDs, nil)

				// 2. The subsequent database call to hydrate the tweets fails.
				storage.EXPECT().
					SelectTweetsByTweetsIDs(gomock.Any(), tweetIDs).
					Return(nil, dbError)
			},
			expectedTweets: nil,
			expectedErr:    dbError,
		},
		{
			name: "Failure - Cache Miss, Fallback Fails when call to SelectFollowersByUserID()",
			setupMocks: func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository) {
				// 1. Cache miss.
				cache.EXPECT().
					LRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]string{}, nil)

				//2. Expect a call to the storage to SelectFollowersByUserID
				storage.EXPECT().SelectFollowersByUserID(gomock.Any(), user1).
					Return(nil, dbError)
			},
			expectedTweets: nil,
			expectedErr:    dbError,
		},
		{
			name: "Failure - Cache Miss, Fallback Fails when call to SelectLastTweetsByUsersID()",
			setupMocks: func(storage *mocks.MockStorageRepo, cache *mocks.MockCacheRepository) {
				// 1. Cache miss.
				cache.EXPECT().
					LRange(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]string{}, nil)

				//2. Expect a call to the storage to SelectFollowersByUserID
				storage.EXPECT().SelectFollowersByUserID(gomock.Any(), user1).
					Return(fallbackFollowers, nil)

				//2. Expect a call to the storage to SelectFollowersByUserID
				storage.EXPECT().SelectLastTweetsByUsersID(gomock.Any(), fallbackFollowers).
					Return(nil, dbError)
			},
			expectedTweets: nil,
			expectedErr:    dbError,
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
			resultTweets, err := service.GetTimeline(context.Background(), user1, limit)

			// Assert
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expectedTweets, resultTweets)
		})
	}
}
