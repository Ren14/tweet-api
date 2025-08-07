package writer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/renzonaitor/tweet-api/internal/domain"
)

const maxTweetLength = 280

type TweetRequest struct {
	Text           string `json:"text"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (h WriterHandler) HandlePublishTweet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
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

	// Validations
	text, err := h.validateMaxLengthText(tweet.Text)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error validating tweet content: %s", err)))
		if err != nil {
			return
		}
		return
	}

	// TODO validate idempotency_key with UUID pattern

	newTweet, err := h.UserService.PublishTweet(r.Context(), domain.Tweet{
		ID:        tweet.IdempotencyKey,
		Text:      text,
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

func (h WriterHandler) validateMaxLengthText(text string) (string, error) {
	// 2. Validate the tweet content.
	// First, trim leading/trailing whitespace to handle empty or space-only tweets.
	trimmedText := strings.TrimSpace(text)
	// Count runes, not bytes, to correctly handle multi-byte characters like emojis.
	runeCount := utf8.RuneCountInString(trimmedText)

	if runeCount == 0 {
		// It's a good practice to return specific validation errors.
		return "", fmt.Errorf("tweet text cannot be empty")
	}

	if runeCount > maxTweetLength {
		return "", fmt.Errorf("tweet exceeds maximum length of %d characters", maxTweetLength)
	}

	// Use the trimmed text for the actual tweet content.
	return trimmedText, nil
}
