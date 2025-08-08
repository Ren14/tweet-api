package writer_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/renzonaitor/tweet-api/cmd/http/handlers/mocks"
	"github.com/renzonaitor/tweet-api/cmd/http/handlers/writer"
	"github.com/renzonaitor/tweet-api/internal/domain"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	xUserID      = "a00ffe35-fc64-45f3-be60-8c824ec0a352"
	userToFollow = "a00ffe35-fc64-45f3-be60-8c824ec0a353"
)

// TestHandleFollowUser uses a table-driven approach with gomock.
func TestHandleFollowUser(t *testing.T) {
	testCases := []struct {
		name                 string
		request              *http.Request
		setupRequest         func(req *http.Request)
		setupMock            func(mock *mocks.MockUserService)
		expectedStatus       int
		expectedBodyContains string
	}{
		{
			name: "Success - 204 No Content",
			request: httptest.NewRequest(
				http.MethodPost,
				"/api/v1/follow",
				strings.NewReader(`{"follow_user_id": "`+userToFollow+`"}`)),
			setupRequest: func(req *http.Request) {
				req.Header.Set("X-User-ID", xUserID)
				req.Header.Set("Content-Type", "application/json")
			},
			setupMock: func(mock *mocks.MockUserService) {
				mock.EXPECT().
					FollowUser(gomock.Any(), domain.FollowUser{
						FollowID:   xUserID,
						FollowedID: userToFollow,
					}).
					Return(nil).
					Times(1)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:                 "Failure - 400 Bad Request for method not allowed",
			request:              httptest.NewRequest(http.MethodGet, "/api/v1/follow", nil),
			setupRequest:         func(req *http.Request) {},
			setupMock:            func(mock *mocks.MockUserService) {},
			expectedStatus:       http.StatusMethodNotAllowed,
			expectedBodyContains: "Method Not Allowed",
		},
		{
			name:                 "Failure - 400 Bad Request for missing user ID header",
			request:              httptest.NewRequest(http.MethodPost, "/api/v1/follow", nil),
			setupRequest:         func(req *http.Request) {},
			setupMock:            func(mock *mocks.MockUserService) {},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "Header X-User-ID is required",
		},
		{
			name: "Failure - 400 Bad Request for error reading body",
			request: httptest.NewRequest(
				http.MethodPost,
				"/api/v1/follow",
				strings.NewReader(`{"follow_user_id":`)),
			setupRequest: func(req *http.Request) {
				req.Header.Set("X-User-ID", xUserID)
			},
			setupMock:            func(mock *mocks.MockUserService) {},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "error unmarshalling body: unexpected end of JSON input",
		},
		{
			name: "Failure - 400 Bad Request for following self",
			request: httptest.NewRequest(
				http.MethodPost,
				"/api/v1/follow",
				strings.NewReader(`{"follow_user_id": "`+xUserID+`"}`)),
			setupRequest: func(req *http.Request) {
				req.Header.Set("X-User-ID", xUserID)
			},
			setupMock:            func(mock *mocks.MockUserService) {},
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "user cannot follow themselves",
		},
		{
			name: "Failure - 500 Internal Server Error from service",
			request: httptest.NewRequest(
				http.MethodPost,
				"/api/v1/follow",
				strings.NewReader(`{"follow_user_id": "`+userToFollow+`"}`)),
			setupRequest: func(req *http.Request) {
				req.Header.Set("X-User-ID", xUserID)
			},
			setupMock: func(mock *mocks.MockUserService) {
				mock.EXPECT().
					FollowUser(gomock.Any(), gomock.Any()).
					Return(errors.New("database constraint violation"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: "error following user",
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

			if tc.setupRequest != nil {
				tc.setupRequest(tc.request)
			}

			// Act
			handler.HandleFollowUser(recorder, tc.request)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			if tc.expectedBodyContains != "" {
				assert.True(t, strings.Contains(recorder.Body.String(), tc.expectedBodyContains))
			}
		})
	}
}
