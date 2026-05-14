// Package metrics provides lightweight in-process counters and gauges
// that track depwatch operational statistics at runtime.
package metrics

import "sync/atomic"

// Metrics holds atomic counters for key depwatch events.
type Metrics struct {
	ChecksTotal      atomic.Int64
	ChangesDetected  atomic.Int64
	NotificationsSent atomic.Int64
	NotificationsFailed atomic.Int64
	StoreErrors      atomic.Int64
}

// New returns an initialised Metrics instance.
func New() *Metrics {
	return &Metrics{}
}

// Snapshot returns a point-in-time copy of all counters as a plain struct
// suitable for serialisation or logging.
func (m *Metrics) Snapshot() Snapshot {
	return Snapshot{
		ChecksTotal:         m.ChecksTotal.Load(),
		ChangesDetected:     m.ChangesDetected.Load(),
		NotificationsSent:   m.NotificationsSent.Load(),
		NotificationsFailed: m.NotificationsFailed.Load(),
		StoreErrors:         m.StoreErrors.Load(),
	}
}

// Snapshot is a value-type copy of Metrics at a point in time.
type Snapshot struct {
	ChecksTotal         int64 `json:"checks_total"`
	ChangesDetected     int64 `json:"changes_detected"`
	NotificationsSent   int64 `json:"notifications_sent"`
	NotificationsFailed int64 `json:"notifications_failed"`
	StoreErrors         int64 `json:"store_errors"`
}
