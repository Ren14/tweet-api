package redis_test

import (
	"strconv"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/renzonaitor/tweet-api/cmd/http/config"
	"github.com/renzonaitor/tweet-api/internal/infraestructure/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepository(t *testing.T) {
	t.Run("Success - Connects to mock Redis", func(t *testing.T) {
		// Arrange
		mockRedis, err := miniredis.Run()
		require.NoError(t, err, "Failed to start mock redis")
		port, err := strconv.Atoi(mockRedis.Port())
		assert.NoError(t, err)
		// t.Cleanup ensures the server is closed when the test finishes.
		t.Cleanup(mockRedis.Close)

		cfg := config.Config{
			Redis: config.Redis{
				Host: mockRedis.Host(),
				Port: port,
			},
		}

		// Act
		repo, err := redis.NewRepository(cfg)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, repo)
		assert.NotNil(t, repo.Client)

		repo.Close()
	})

	t.Run("Failure - Fails to connect", func(t *testing.T) {
		// Arrange
		// Create a config that points to a non-existent server.
		// Using a port of 1 ensures nothing is running there.
		cfg := config.Config{
			Redis: config.Redis{
				Host: "127.0.0.1",
				Port: 1,
			},
		}

		// Act
		repo, err := redis.NewRepository(cfg)

		// Assert
		require.Error(t, err)
		assert.Nil(t, repo)
		assert.Contains(t, err.Error(), "failed to connect to Redis")
	})
}
