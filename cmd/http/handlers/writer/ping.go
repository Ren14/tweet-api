package writer

import (
	"net/http"
)

func (h WriterHandler) Ping(w http.ResponseWriter, r *http.Request) {
	// TODO validate HTTP method

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
