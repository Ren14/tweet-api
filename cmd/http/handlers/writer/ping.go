package writer

import (
	"net/http"
)

func (h WriterHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// If not, respond with a 405 Method Not Allowed error.
		w.Header().Set("Allow", http.MethodGet) // Let the client know which method is allowed.
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return // Stop further execution.
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		return
	}
}
