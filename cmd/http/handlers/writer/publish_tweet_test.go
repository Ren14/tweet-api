package writer_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/renzonaitor/tweet-api/cmd/http/handlers/mocks"
	"github.com/renzonaitor/tweet-api/cmd/http/handlers/writer"
	"github.com/renzonaitor/tweet-api/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestHandlePublishTweet uses a table-driven approach with gomock.
func TestHandlePublishTweet(t *testing.T) {
	// Define reusable constants for test clarity
	const testUserID = "a00ffe35-fc64-45f3-be60-8c824ec0a352"
	idempotencyKey := uuid.NewString()

	// Define a successful tweet response object to be returned by the mock
	mockTweetResponse := domain.Tweet{
		ID:        idempotencyKey,
		Text:      "This is a valid tweet!",
		UserID:    testUserID,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	testCases := []struct {
		name                 string
		body                 string
		setupRequest         func(req *http.Request)
		setupMock            func(mock *mocks.MockUserService)
		expectedStatus       int
		expectedBodyContains string
		expectedJSONResponse *domain.Tweet
	}{
		{
			name: "Success - 200 OK",
			body: `{"text": "This is a valid tweet!", "idempotency_key": "` + idempotencyKey + `"}`,
			setupRequest: func(req *http.Request) {
				req.Header.Set("X-User-ID", testUserID)
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock: func(mock *mocks.MockUserService) {
				mock.EXPECT().
					PublishTweet(gomock.Any(), gomock.Any()).
					Return(mockTweetResponse, nil).
					Times(1)
			},
			expectedStatus:       http.StatusOK,
			expectedJSONResponse: &mockTweetResponse,
		},
		{
			name:                 "Failure - 405 Method Not Allowed",
			body:                 "",
			setupRequest:         func(req *http.Request) {},
			setupMock:            func(mock *mocks.MockUserService) {},
			expectedStatus:       http.StatusMethodNotAllowed,
			expectedBodyContains: "Method Not Allowed",
		},
		{
			name:                 "Failure - 400 Bad Request for missing user ID header",
			body:                 `{"text": "some text"}`,
			setupRequest:         func(req *http.Request) { req.Header.Set("Content-Type", "application/json") },
			setupMock:            func(mock *mocks.MockUserService) {},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "Header X-User-ID is required",
		},
		{
			name:                 "Failure - 400 Bad Request for malformed JSON",
			body:                 `{"text": "some text"`,
			setupRequest:         func(req *http.Request) { req.Header.Set("X-User-ID", testUserID) },
			setupMock:            func(mock *mocks.MockUserService) {},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "error unmarshalling body",
		},
		{
			name:                 "Failure - 400 Bad Request for empty tweet text",
			body:                 `{"text": "   ", "idempotency_key": "` + idempotencyKey + `"}`,
			setupRequest:         func(req *http.Request) { req.Header.Set("X-User-ID", testUserID) },
			setupMock:            func(mock *mocks.MockUserService) {},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "tweet text cannot be empty",
		},
		{
			name:                 "Failure - 400 Bad Request for tweet text exceeding max length",
			body:                 `{"text": "` + strings.Repeat("a", 281) + `"}`,
			setupRequest:         func(req *http.Request) { req.Header.Set("X-User-ID", testUserID) },
			setupMock:            func(mock *mocks.MockUserService) {},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "tweet exceeds maximum length of 280 characters",
		},
		{
			name: "Failure - 500 Internal Server Error from service",
			body: `{"text": "a valid tweet", "idempotency_key": "` + idempotencyKey + `"}`,
			setupRequest: func(req *http.Request) {
				req.Header.Set("X-User-ID", testUserID)
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock: func(mock *mocks.MockUserService) {
				mock.EXPECT().
					PublishTweet(gomock.Any(), gomock.Any()).
					Return(domain.Tweet{}, errors.New("idempotency key already exists"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: "error publishing tweet",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockService := mocks.NewMockUserService(ctrl)
			tc.setupMock(mockService)

			handler := writer.NewHandler(mockService)
			recorder := httptest.NewRecorder()

			// Create the request
			var request *http.Request
			if tc.name == "Failure - 405 Method Not Allowed" {
				request = httptest.NewRequest(http.MethodGet, "/api/v1/tweets", nil)
			} else {
				request = httptest.NewRequest(http.MethodPost, "/api/v1/tweets", strings.NewReader(tc.body))
			}

			if tc.setupRequest != nil {
				tc.setupRequest(request)
			}

			// Act
			handler.HandlePublishTweet(recorder, request)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)

			if tc.expectedBodyContains != "" {
				assert.Contains(t, recorder.Body.String(), tc.expectedBodyContains)
			}

			if tc.expectedJSONResponse != nil {
				expectedJSON, err := json.Marshal(tc.expectedJSONResponse)
				require.NoError(t, err)
				assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
				assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
			}
		})
	}
}
