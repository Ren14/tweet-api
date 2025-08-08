package redis_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/renzonaitor/tweet-api/cmd/http/config"
	"github.com/renzonaitor/tweet-api/internal/infraestructure/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepo is a helper function to create a mock redis server and a repository instance for tests.
func setupTestRepo(t *testing.T) (*redis.Repository, *miniredis.Miniredis) {
	t.Helper()

	mockRedis, err := miniredis.Run()
	require.NoError(t, err, "Failed to start mock redis")
	port, err := strconv.Atoi(mockRedis.Port())
	assert.NoError(t, err)

	cfg := config.Config{
		Redis: config.Redis{
			Host: mockRedis.Host(),
			Port: port,
		},
	}

	// Act
	repo, err := redis.NewRepository(cfg)
	assert.NoError(t, err)

	return repo, mockRedis
}

func TestLRange(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - retrieves a full list", func(t *testing.T) {
		// Arrange
		repo, mockRedis := setupTestRepo(t)
		t.Cleanup(mockRedis.Close)
		userID, err := uuid.NewUUID()
		assert.Nil(t, err)
		tweet1, err := uuid.NewUUID()
		assert.Nil(t, err)
		tweet2, err := uuid.NewUUID()
		assert.Nil(t, err)
		tweet3, err := uuid.NewUUID()
		assert.Nil(t, err)

		listKey := fmt.Sprintf("tweet_list_%s", userID)
		expectedValues := []string{tweet3.String(), tweet2.String(), tweet1.String()}

		// Pre-populate the mock redis with data for the test.
		// LPush adds to the head, so we push in reverse order.
		err = repo.LPush(context.Background(), listKey, tweet1.String())
		assert.NoError(t, err)
		err = repo.LPush(context.Background(), listKey, tweet2.String())
		assert.NoError(t, err)
		err = repo.LPush(context.Background(), listKey, tweet3.String())
		assert.NoError(t, err)

		// Act
		result, err := repo.LRange(ctx, listKey, 0, 2)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedValues, result)
	})

	t.Run("Success - retrieves a partial list", func(t *testing.T) {
		// Arrange
		repo, mockRedis := setupTestRepo(t)
		t.Cleanup(mockRedis.Close)
		userID, err := uuid.NewUUID()
		assert.Nil(t, err)
		tweet1, err := uuid.NewUUID()
		assert.Nil(t, err)
		tweet2, err := uuid.NewUUID()
		assert.Nil(t, err)
		tweet3, err := uuid.NewUUID()
		assert.Nil(t, err)

		listKey := fmt.Sprintf("tweet_list_%s", userID)
		expectedValues := []string{tweet3.String(), tweet2.String()}

		// Pre-populate the mock redis with data for the test.
		// LPush adds to the head, so we push in reverse order.
		err = repo.LPush(context.Background(), listKey, tweet1.String())
		assert.NoError(t, err)
		err = repo.LPush(context.Background(), listKey, tweet2.String())
		assert.NoError(t, err)
		err = repo.LPush(context.Background(), listKey, tweet3.String())
		assert.NoError(t, err)

		// Act
		result, err := repo.LRange(ctx, listKey, 0, 1)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedValues, result)
	})

	t.Run("Success - key does not exist", func(t *testing.T) {
		// Arrange
		repo, mockRedis := setupTestRepo(t)
		t.Cleanup(mockRedis.Close)

		// Act
		result, err := repo.LRange(ctx, "non-existent-key", 0, -1)

		// Assert
		require.NoError(t, err)
		// Redis LRange on a non-existent key returns an empty list, not an error.
		assert.Empty(t, result)
		assert.NotNil(t, result, "Should return an empty slice, not a nil slice")
	})

	t.Run("Failure - connection error", func(t *testing.T) {
		// Arrange
		repo, mockRedis := setupTestRepo(t)
		// Close the server immediately to simulate a connection failure.
		mockRedis.Close()

		// Act
		result, err := repo.LRange(ctx, "any-key", 0, -1)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to LRANGE")
	})
}

func TestLRange2(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	tweet1 := uuid.NewString()
	tweet2 := uuid.NewString()
	tweet3 := uuid.NewString()

	testCases := []struct {
		name           string
		listKey        string
		initialData    []string
		start          int64
		stop           int64
		setup          func(mr *miniredis.Miniredis)
		expectedResult []string
		expectError    bool
		errorContains  string
	}{
		{
			name:           "Success - retrieves a full list",
			listKey:        fmt.Sprintf("timeline:%s", userID),
			initialData:    []string{tweet1, tweet2, tweet3},
			start:          0,
			stop:           2,
			expectedResult: []string{tweet3, tweet2, tweet1},
			expectError:    false,
		},
		{
			name:           "Success - retrieves a partial list",
			listKey:        fmt.Sprintf("timeline:%s", userID),
			initialData:    []string{tweet1, tweet2, tweet3},
			start:          0,
			stop:           1,
			expectedResult: []string{tweet3, tweet2},
			expectError:    false,
		},
		{
			name:           "Success - key does not exist",
			listKey:        "non-existent-key",
			initialData:    nil, // No data to pre-populate
			start:          0,
			stop:           1,
			expectedResult: []string{}, // Expect an empty slice, not nil
			expectError:    false,
		},
		{
			name:    "Failure - connection error",
			listKey: "any-key",
			start:   0,
			stop:    1,
			setup: func(mr *miniredis.Miniredis) {
				// Simulate a connection failure by closing the server.
				mr.Close()
			},
			expectedResult: nil,
			expectError:    true,
			errorContains:  "failed to LRANGE",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo, mockRedis := setupTestRepo(t)
			t.Cleanup(mockRedis.Close)

			// Pre-populate the mock redis with data for the test.
			if len(tc.initialData) > 0 {
				// LPush adds to the head, so we push in reverse order of how we want to read it.
				for i := 0; i < len(tc.initialData); i++ {
					err := repo.LPush(ctx, tc.listKey, tc.initialData[i])
					require.NoError(t, err)
				}
			}

			if tc.setup != nil {
				tc.setup(mockRedis)
			}

			// Act
			result, err := repo.LRange(ctx, tc.listKey, tc.start, tc.stop)

			// Assert
			if tc.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
