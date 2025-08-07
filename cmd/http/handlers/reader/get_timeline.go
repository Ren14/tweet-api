package reader

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const defaultPaginationLimit = 10

func (h *ReaderHandler) HandleGetTimeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
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

	limit, err := parseLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("error parsing limit: %s", err)))
		if err != nil {
			return
		}
		return
	}

	// nextCursor := r.URL.Query().Get("next_cursor") // TODO implement pagination

	timeline, err := h.Timeline.GetTimeline(r.Context(), userID, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(fmt.Sprintf("error getting timeline: %s", err)))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	timelineResponse, err := json.Marshal(timeline)
	if err != nil {
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
	}

	_, err = w.Write(timelineResponse)
	if err != nil {
		return
	}
}

func parseLimit(request *http.Request) (int, error) {
	limitResponse := defaultPaginationLimit
	limitStr := request.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, fmt.Errorf("limit must be an integer: %w", err)
		}
		if limit <= 0 {
			return 0, errors.New("limit must be a positive number")
		}
		limitResponse = limit
	}

	return limitResponse, nil
}
