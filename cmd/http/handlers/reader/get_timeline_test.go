package reader_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/renzonaitor/tweet-api/cmd/http/handlers/mocks"
	"github.com/renzonaitor/tweet-api/cmd/http/handlers/reader"
	"github.com/renzonaitor/tweet-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestHandleGetTimeline uses a table-driven approach with gomock.
func TestHandleGetTimeline(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	mockTweets := []domain.Tweet{
		{ID: "a00ffe35-fc64-45f3-be60-8c824ec0a353", UserID: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", Text: "Hello World", CreatedAt: now},
		{ID: "a00ffe35-fc64-45f3-be60-8c824ec0a352", UserID: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12", Text: "Testing handlers!", CreatedAt: now},
	}

	testCases := []struct {
		name                 string
		setupMock            func(mock *mocks.MockTimelineService)
		request              *http.Request
		setupRequest         func(req *http.Request)
		expectedStatus       int
		expectedBodyContains string
		expectedJSONResponse []domain.Tweet
	}{
		{
			name: "Success - 200 OK with default limit",
			setupMock: func(mock *mocks.MockTimelineService) {
				mock.EXPECT().
					GetTimeline(gomock.Any(), "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", 10).
					Return(mockTweets, nil).
					Times(1)
			},
			request:              httptest.NewRequest(http.MethodGet, "/timeline", nil),
			setupRequest:         func(req *http.Request) { req.Header.Set("X-User-ID", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11") },
			expectedStatus:       http.StatusOK,
			expectedJSONResponse: mockTweets,
		},
		{
			name: "Success - 200 OK with custom limit",
			setupMock: func(mock *mocks.MockTimelineService) {
				mock.EXPECT().
					GetTimeline(gomock.Any(), "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", 5).
					Return(mockTweets, nil)
			},
			request:              httptest.NewRequest(http.MethodGet, "/timeline?limit=5", nil),
			setupRequest:         func(req *http.Request) { req.Header.Set("X-User-ID", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11") },
			expectedStatus:       http.StatusOK,
			expectedJSONResponse: mockTweets,
		},
		{
			name:                 "Failure - 400 Bad Request for missing user ID",
			setupMock:            func(mock *mocks.MockTimelineService) {}, // No calls to the mock are expected
			request:              httptest.NewRequest(http.MethodGet, "/timeline", nil),
			setupRequest:         func(req *http.Request) {},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "Header X-User-ID is required",
		},
		{
			name:                 "Failure - 400 Bad Request for missing limit",
			setupMock:            func(mock *mocks.MockTimelineService) {}, // No calls to the mock are expected
			request:              httptest.NewRequest(http.MethodGet, "/timeline?limit=-1", nil),
			setupRequest:         func(req *http.Request) { req.Header.Set("X-User-ID", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11") },
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "error parsing limit: limit must be a positive number",
		},
		{
			name:                 "Failure - 400 Bad Request for method not allowed",
			setupMock:            func(mock *mocks.MockTimelineService) {}, // No calls to the mock are expected
			request:              httptest.NewRequest(http.MethodPost, "/timeline", nil),
			setupRequest:         func(req *http.Request) { req.Header.Set("X-User-ID", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11") },
			expectedStatus:       http.StatusMethodNotAllowed,
			expectedBodyContains: "Method Not Allowed\n",
		},
		{
			name: "Failure - 500 Internal Server Error from service",
			setupMock: func(mock *mocks.MockTimelineService) {
				mock.EXPECT().
					GetTimeline(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database is down"))
			},
			request:              httptest.NewRequest(http.MethodGet, "/timeline", nil),
			setupRequest:         func(req *http.Request) { req.Header.Set("X-User-ID", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11") },
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: "error getting timeline: database is down",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockService := mocks.NewMockTimelineService(ctrl)
			tc.setupMock(mockService)

			handler := reader.NewHandler(mockService)
			recorder := httptest.NewRecorder()

			if tc.setupRequest != nil {
				tc.setupRequest(tc.request)
			}

			// Act
			handler.HandleGetTimeline(recorder, tc.request)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)

			if tc.expectedBodyContains != "" {
				assert.Equal(t, recorder.Body.String(), tc.expectedBodyContains)
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
