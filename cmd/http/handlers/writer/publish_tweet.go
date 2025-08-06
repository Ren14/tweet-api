package writer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

type TweetRequest struct {
	Text           string `json:"text"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (h WriterHandler) HandlePublishTweet(w http.ResponseWriter, r *http.Request) {
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

	var tweet TweetRequest
	if err = json.Unmarshal(bytes, &tweet); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error unmarshalling body: %s", err)))
		if err != nil {
			return
		}
		return
	}

	newTweet, err := h.UserService.PublishTweet(r.Context(), domain.Tweet{
		ID:        tweet.IdempotencyKey,
		Text:      tweet.Text,
		UserID:    userID,
		CreatedAt: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(fmt.Sprintf("error publishing tweet: %s", err)))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	tweetResponse, err := json.Marshal(newTweet)
	if err != nil {
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
	}

	_, err = w.Write(tweetResponse)
	if err != nil {
		return
	}
}
