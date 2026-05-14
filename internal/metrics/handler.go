// Package metrics provides counters and snapshot reporting for depwatch.
package metrics

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.HandlerFunc that serves a JSON snapshot of m.
//
// GET /metrics
//
//	{
//	  "checks_total": 42,
//	  "updates_detected": 3,
//	  "notifications_sent": 3,
//	  "notification_errors": 0
//	}
func (m *Metrics) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap := m.Snapshot()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(snap); err != nil {
			http.Error(w, "failed to encode metrics", http.StatusInternalServerError)
		}
	}
}
