package omikuji

import (
	"encoding/json"
	"net/http"
)

// Handler is an HTTP handler that serves omikuji fortune responses.
type Handler struct {
	clock Clock
}

// NewHandler creates a new Handler with the specified clock function.
func NewHandler(clock Clock) *Handler {
	return &Handler{clock: clock}
}

// ServeHTTP handles HTTP requests and returns a random fortune in JSON format.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fortune := GenerateFortune(h.clock)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(fortune); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
