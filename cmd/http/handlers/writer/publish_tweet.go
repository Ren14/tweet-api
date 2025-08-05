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

	var tweet TweetRequest
	if err := json.Unmarshal(bytes, &tweet); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error unmarshalling body: %s", err)))
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
		w.Write([]byte(fmt.Sprintf("error publishing tweet: %s", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	tweetResponse, err := json.Marshal(newTweet)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(tweetResponse)
}
