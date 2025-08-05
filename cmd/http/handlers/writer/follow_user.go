package writer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

type FollowUserRequest struct {
	FollowUserID string `json:"follow_user_id"`
}

func (h *WriterHandler) HandleFollowUser(w http.ResponseWriter, r *http.Request) {
	// TODO validate HTTP method

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Header X-User-ID is required"))
		return
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error reading body: %s", err)))
		return
	}

	var follow FollowUserRequest
	if err := json.Unmarshal(bytes, &follow); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error unmarshalling body: %s", err)))
		return
	}

	err = h.UserService.FollowUser(r.Context(), domain.FollowUser{
		FollowID:   userID,
		FollowedID: follow.FollowUserID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error following user: %s", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
