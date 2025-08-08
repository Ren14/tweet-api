package postgres_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/renzonaitor/tweet-api/cmd/http/config"
	"github.com/renzonaitor/tweet-api/internal/infraestructure/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewRepositoryWithDB uses a table-driven approach to test the repository constructor.
func TestNewRepositoryWithDB(t *testing.T) {
	cfg := config.Config{
		Postgres: config.Postgres{
			MaxOpenConnection: 10,
			MaxIdleConnection: 5,
		},
	}

	testCases := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		expectError   bool
		errorContains string
	}{
		{
			name: "Success - connects and pings database",
			setupMock: func(mock sqlmock.Sqlmock) {
				// We expect a Ping to be made. The mock will return no error.
				mock.ExpectPing()
			},
			expectError: false,
		},
		{
			name: "Failure - fails to ping database",
			setupMock: func(mock sqlmock.Sqlmock) {
				// We expect a Ping, but tell the mock to return an error.
				pingError := errors.New("database is offline")
				mock.ExpectPing().WillReturnError(pingError)

				// We also expect the Close call that happens on ping failure.
				mock.ExpectClose()
			},
			expectError:   true,
			errorContains: "failed to ping database",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			require.NoError(t, err)
			defer db.Close()

			if tc.setupMock != nil {
				tc.setupMock(mock)
			}

			repo, err := postgres.NewRepositoryWithDB(db, cfg)

			// Assert
			if tc.expectError {
				require.Error(t, err)
				assert.Nil(t, repo)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, repo)
			}
		})
	}
}
