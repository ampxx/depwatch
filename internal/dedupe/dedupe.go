// Package dedupe provides a deduplication layer that suppresses repeated
// alerts for the same module/version pair within a configurable window.
package dedupe

import (
	"sync"
	"time"
)

// entry records when an alert was last sent for a given key.
type entry struct {
	version   string
	sentAt    time.Time
}

// Deduplicator tracks recently-sent alerts and suppresses duplicates.
type Deduplicator struct {
	mu      sync.Mutex
	window  time.Duration
	records map[string]entry
}

// New creates a Deduplicator that suppresses repeated alerts within window.
func New(window time.Duration) *Deduplicator {
	return &Deduplicator{
		window:  window,
		records: make(map[string]entry),
	}
}

// IsDuplicate reports whether an alert for module/version was already sent
// within the deduplication window. If it is not a duplicate the record is
// updated so subsequent calls within the window return true.
func (d *Deduplicator) IsDuplicate(module, version string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	key := module + "@" + version
	if e, ok := d.records[key]; ok {
		if time.Since(e.sentAt) < d.window {
			return true
		}
	}

	d.records[key] = entry{version: version, sentAt: time.Now()}
	return false
}

// Purge removes all entries whose window has expired. It is safe to call
// periodically to prevent unbounded memory growth.
func (d *Deduplicator) Purge() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for key, e := range d.records {
		if time.Since(e.sentAt) >= d.window {
			delete(d.records, key)
		}
	}
}

// Len returns the number of active (non-expired) records.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	count := 0
	for _, e := range d.records {
		if time.Since(e.sentAt) < d.window {
			count++
		}
	}
	return count
}
