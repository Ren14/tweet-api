package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/renzonaitor/tweet-api/internal/domain"
	"github.com/renzonaitor/tweet-api/internal/service/user"
	"github.com/renzonaitor/tweet-api/internal/service/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFollowUser(t *testing.T) {
	followInput := domain.FollowUser{
		FollowID:   uuid.NewString(),
		FollowedID: uuid.NewString(),
	}

	// Define a reusable database error
	dbError := errors.New("database constraint violation: user already follows this user")

	testCases := []struct {
		name        string
		input       domain.FollowUser
		setupMock   func(storage *mocks.MockStorageRepo)
		expectedErr error
	}{
		{
			name:  "Success - Create Follow Relation",
			input: followInput,
			setupMock: func(storage *mocks.MockStorageRepo) {
				storage.EXPECT().
					CreateRelation(gomock.Any(), followInput).
					Return(nil).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name:  "Failure - Error from storage layer",
			input: followInput,
			setupMock: func(storage *mocks.MockStorageRepo) {
				storage.EXPECT().
					CreateRelation(gomock.Any(), followInput).
					Return(dbError)
			},
			expectedErr: dbError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockStorage := mocks.NewMockStorageRepo(ctrl)

			if tc.setupMock != nil {
				tc.setupMock(mockStorage)
			}

			service := user.NewService(mockStorage, nil)

			// Act
			err := service.FollowUser(context.Background(), tc.input)

			// Assert
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
