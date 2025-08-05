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
	// TODO validate HTTP method

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Header X-User-ID is required"))
		return
	}

	limit, err := parseLimit(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error parsing limit: %s", err)))
	}

	// nextCursor := r.URL.Query().Get("next_cursor") // TODO implement pagination

	timeline, err := h.Timeline.GetTimeline(r.Context(), userID, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error getting timeline: %s", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	timelineResponse, err := json.Marshal(timeline)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(timelineResponse)
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
