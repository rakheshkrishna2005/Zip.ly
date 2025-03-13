package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// HealthCheck returns the status of the service
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status    string        `json:"status"`
		Uptime    string        `json:"uptime"`
		Timestamp time.Time     `json:"timestamp"`
		Duration  time.Duration `json:"duration_ms"`
	}{
		Status:    "ok",
		Uptime:    time.Since(h.startTime).String(),
		Timestamp: time.Now(),
		Duration:  time.Since(time.Now().Add(-1 * time.Millisecond)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}