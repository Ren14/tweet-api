package user_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/renzonaitor/tweet-api/internal/domain"
	"github.com/renzonaitor/tweet-api/internal/service/user"
	"github.com/renzonaitor/tweet-api/internal/service/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPublishTweet(t *testing.T) {
	inputTweet := domain.Tweet{
		ID:     uuid.NewString(),
		UserID: uuid.NewString(),
		Text:   "This is a test tweet!",
	}

	dbError := errors.New("database connection lost")

	testCases := []struct {
		name          string
		input         domain.Tweet
		setupMocks    func(storage *mocks.MockStorageRepo, timeline *mocks.MockTimelineUpdater, wg *sync.WaitGroup)
		expectedTweet domain.Tweet
		expectedErr   error
	}{
		{
			name:  "Success - New Tweet",
			input: inputTweet,
			setupMocks: func(storage *mocks.MockStorageRepo, timeline *mocks.MockTimelineUpdater, wg *sync.WaitGroup) {
				// 1. Expect a call to check for the tweet, and it's not found.
				storage.EXPECT().SelectTweetByID(gomock.Any(), inputTweet.ID).Return(nil, nil)

				// 2. Expect a call to create the tweet, which succeeds.
				storage.EXPECT().CreateTweet(gomock.Any(), inputTweet).Return(inputTweet, nil)

				// 3. Expect the async timeline update. We use a WaitGroup to test this.
				wg.Add(1)
				timeline.EXPECT().
					UpdateTimeline(gomock.Any(), inputTweet.UserID, inputTweet.ID).
					Do(func(ctx context.Context, authorID, tweetID string) {
						wg.Done() // Signal that the mock was called
					})
			},
			expectedTweet: inputTweet,
			expectedErr:   nil,
		},
		{
			name:  "Success - Idempotency Hit",
			input: inputTweet,
			setupMocks: func(storage *mocks.MockStorageRepo, timeline *mocks.MockTimelineUpdater, wg *sync.WaitGroup) {
				storage.EXPECT().SelectTweetByID(gomock.Any(), inputTweet.ID).Return(&inputTweet, nil)
				// No other calls to storage or timeline are expected.
			},
			expectedTweet: inputTweet,
			expectedErr:   nil,
		},
		{
			name:  "Failure - Error checking for existing tweet",
			input: inputTweet,
			setupMocks: func(storage *mocks.MockStorageRepo, timeline *mocks.MockTimelineUpdater, wg *sync.WaitGroup) {
				storage.EXPECT().SelectTweetByID(gomock.Any(), inputTweet.ID).Return(nil, dbError)
			},
			expectedTweet: domain.Tweet{},
			expectedErr:   dbError,
		},
		{
			name:  "Failure - Error creating tweet",
			input: inputTweet,
			setupMocks: func(storage *mocks.MockStorageRepo, timeline *mocks.MockTimelineUpdater, wg *sync.WaitGroup) {
				storage.EXPECT().SelectTweetByID(gomock.Any(), inputTweet.ID).Return(nil, nil)
				storage.EXPECT().CreateTweet(gomock.Any(), inputTweet).Return(domain.Tweet{}, dbError)
			},
			expectedTweet: domain.Tweet{},
			expectedErr:   dbError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockStorage := mocks.NewMockStorageRepo(ctrl)
			mockTimeline := mocks.NewMockTimelineUpdater(ctrl)
			wg := &sync.WaitGroup{}

			if tc.setupMocks != nil {
				tc.setupMocks(mockStorage, mockTimeline, wg)
			}

			service := user.NewService(mockStorage, mockTimeline)

			// Act
			resultTweet, err := service.PublishTweet(context.Background(), tc.input)

			// Assert
			// Wait for the goroutine to finish (if one was expected)
			wg.Wait()

			// Check the error
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			// Check the returned tweet
			assert.Equal(t, tc.expectedTweet, resultTweet)
		})
	}
}
