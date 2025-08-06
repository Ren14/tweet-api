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
	if r.Method != http.MethodPost {
		// If not, respond with a 405 Method Not Allowed error.
		w.Header().Set("Allow", http.MethodPost) // Let the client know which method is allowed.
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return // Stop further execution.
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Header X-User-ID is required"))
		if err != nil {
			return
		}
		return
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error reading body: %s", err)))
		if err != nil {
			return
		}
		return
	}

	var follow FollowUserRequest
	if err := json.Unmarshal(bytes, &follow); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error unmarshalling body: %s", err)))
		if err != nil {
			return
		}
		return
	}

	err = h.UserService.FollowUser(r.Context(), domain.FollowUser{
		FollowID:   userID,
		FollowedID: follow.FollowUserID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(fmt.Sprintf("error following user: %s", err)))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
