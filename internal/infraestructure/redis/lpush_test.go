package redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLPush(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	tweet1 := uuid.NewString()
	tweet2 := uuid.NewString()
	tweet3 := uuid.NewString()
	timelineKey := fmt.Sprintf("timeline:%s", userID)

	testCases := []struct {
		name              string
		timelineKey       string
		initialData       []interface{}
		valuesToPush      []interface{}
		setup             func(mr *miniredis.Miniredis)
		expectedListState []string
		expectError       bool
		errorContains     string
	}{
		{
			name:              "Success - push single value to new list",
			timelineKey:       timelineKey,
			valuesToPush:      []interface{}{tweet1},
			expectedListState: []string{tweet1},
			expectError:       false,
		},
		{
			name:              "Success - push multiple values to new list",
			timelineKey:       timelineKey,
			valuesToPush:      []interface{}{tweet2, tweet3},
			expectedListState: []string{tweet3, tweet2}, // LPush prepends, so the last value pushed is the first in the list.
			expectError:       false,
		},
		{
			name:              "Success - push to an existing list",
			timelineKey:       timelineKey,
			initialData:       []interface{}{tweet1},
			valuesToPush:      []interface{}{tweet2},
			expectedListState: []string{tweet2, tweet1},
			expectError:       false,
		},
		{
			name:        "Failure - connection error",
			timelineKey: "any-key",
			setup: func(mr *miniredis.Miniredis) {
				// Simulate a connection failure by closing the server.
				mr.Close()
			},
			valuesToPush:  []interface{}{"tweet-f"},
			expectError:   true,
			errorContains: "failed to LPUSH",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo, mockRedis := setupTestRepo(t)
			t.Cleanup(mockRedis.Close)

			// Pre-populate the mock redis with initial data if needed.
			if len(tc.initialData) > 0 {
				err := repo.LPush(ctx, tc.timelineKey, tc.initialData...)
				require.NoError(t, err)
			}

			if tc.setup != nil {
				tc.setup(mockRedis)
			}

			// Act
			err := repo.LPush(ctx, tc.timelineKey, tc.valuesToPush...)

			// Assert
			if tc.expectError {
				require.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				require.NoError(t, err)
				// Verify the final state of the list in Redis to confirm the push was successful.
				actualListState, err := repo.LRange(ctx, tc.timelineKey, 0, -1)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedListState, actualListState)
			}
		})
	}
}
