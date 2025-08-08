package user_test

import (
	"testing"

	"github.com/renzonaitor/tweet-api/internal/service/user"
	"github.com/renzonaitor/tweet-api/internal/service/user/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestNewService verifies that the service constructor correctly initializes
// the service with its dependencies.
func TestNewService(t *testing.T) {
	// Arrange: Set up the gomock controller and create mock dependencies.
	ctrl := gomock.NewController(t)

	// Create mock instances using the auto-generated constructors.
	mockStorage := mocks.NewMockStorageRepo(ctrl)
	mockTimeline := mocks.NewMockTimelineUpdater(ctrl)

	// Act: Call the constructor function that we are testing.
	service := user.NewService(mockStorage, mockTimeline)

	// Assert: Verify the outcome.
	// 1. Ensure the service object was actually created.
	assert.NotNil(t, service)

	// 2. Ensure the dependencies were assigned to the correct fields.
	// This confirms that the service holds the dependencies it needs to operate.
	assert.Equal(t, mockStorage, service.Storage, "Storage should be the provided mock instance")
	assert.Equal(t, mockTimeline, service.Timeline, "TimelineUpdater should be the provided mock instance")
}
