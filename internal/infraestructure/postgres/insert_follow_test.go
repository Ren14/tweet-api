package postgres_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/renzonaitor/tweet-api/cmd/http/config"
	"github.com/renzonaitor/tweet-api/internal/domain"
	"github.com/renzonaitor/tweet-api/internal/infraestructure/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupRepoWithMock is a test helper to create a repository instance with a mock DB.
// It handles the boilerplate of satisfying the constructor's ping expectation.
func setupRepoWithMock(t *testing.T) (*postgres.Repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	// Satisfy the ping from the NewRepositoryWithDB constructor
	mock.ExpectPing()

	cfg := config.Config{
		Postgres: config.Postgres{
			MaxOpenConnection: 10,
			MaxIdleConnection: 5,
		},
	}
	repo, err := postgres.NewRepositoryWithDB(db, cfg)
	require.NoError(t, err)

	return repo, mock
}

func TestCreateRelation(t *testing.T) {
	ctx := context.Background()
	followInput := domain.FollowUser{
		FollowID:   uuid.NewString(),
		FollowedID: uuid.NewString(),
	}

	// The SQL query we expect to be executed. Using regexp.QuoteMeta makes it
	// robust against whitespace changes.
	expectedQuery := regexp.QuoteMeta(`
		INSERT INTO follows (follower_id, following_id)
		VALUES ($1, $2)
	`)

	// A simulated unique constraint violation error from Postgres.
	uniqueConstraintErr := errors.New(`pq: duplicate key value violates unique constraint "follows_pkey"`)

	testCases := []struct {
		name          string
		input         domain.FollowUser
		setupMock     func(mock sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name:  "Success - creates follow relation",
			input: followInput,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(followInput.FollowID, followInput.FollowedID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name:  "Failure - database error on exec",
			input: followInput,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(followInput.FollowID, followInput.FollowedID).
					WillReturnError(errors.New("database connection lost"))
			},
			expectError:   true,
			errorContains: "database connection lost",
		},
		{
			name:  "Failure - unique constraint violation",
			input: followInput,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(followInput.FollowID, followInput.FollowedID).
					WillReturnError(uniqueConstraintErr)
			},
			expectError:   true,
			errorContains: `violates unique constraint "follows_pkey"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			repo, mock := setupRepoWithMock(t)
			tc.setupMock(mock)

			// Act
			err := repo.CreateRelation(ctx, tc.input)

			// Assert
			if tc.expectError {
				require.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
