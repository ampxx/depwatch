// Package healthcheck provides an HTTP handler for liveness and readiness probes.
package healthcheck

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// Status holds the current health state of the daemon.
type Status struct {
	OK        bool      `json:"ok"`
	Uptime    string    `json:"uptime"`
	CheckedAt time.Time `json:"checked_at"`
}

// Handler is an HTTP handler that reports service health.
type Handler struct {
	start   time.Time
	healthy atomic.Bool
}

// New creates a new Handler. The service is considered healthy by default.
func New() *Handler {
	h := &Handler{start: time.Now()}
	h.healthy.Store(true)
	return h
}

// SetHealthy marks the service as healthy or unhealthy.
func (h *Handler) SetHealthy(ok bool) {
	h.healthy.Store(ok)
}

// ServeHTTP writes a JSON health status response.
// It returns HTTP 200 when healthy and HTTP 503 when unhealthy.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ok := h.healthy.Load()
	status := Status{
		OK:        ok,
		Uptime:    time.Since(h.start).Round(time.Second).String(),
		CheckedAt: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	if !ok {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	_ = json.NewEncoder(w).Encode(status)
}
