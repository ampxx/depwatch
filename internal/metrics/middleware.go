package metrics

import (
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Middleware returns an http.Handler that records request metrics.
// It increments ChecksTotal on every request and tracks latency via
// the provided Metrics instance. Non-2xx responses increment AlertsTotal
// as a proxy for error accounting.
func Middleware(m *Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		m.mu.Lock()
		m.data.ChecksTotal++
		if m.data.LastCheckDuration == 0 || duration > m.data.LastCheckDuration {
			m.data.LastCheckDuration = duration
		}
		if rw.statusCode < 200 || rw.statusCode >= 300 {
			m.data.AlertsTotal++
		}
		m.mu.Unlock()
	})
}
