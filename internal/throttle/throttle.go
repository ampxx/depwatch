// Package throttle provides a per-module notification throttle that
// suppresses repeated alerts for the same module within a configurable
// cooldown window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks the last notification time per module and suppresses
// alerts that occur within the cooldown window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSent map[string]time.Time
	now      func() time.Time
}

// New returns a Throttle with the given cooldown duration.
func New(cooldown time.Duration) *Throttle {
	return &Throttle{
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if a notification for the given module should be
// sent, i.e. no notification has been sent within the cooldown window.
// Calling Allow with a permitted module records the current time so
// subsequent calls within the window are suppressed.
func (t *Throttle) Allow(module string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.lastSent[module]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.lastSent[module] = now
	return true
}

// Reset clears the recorded send time for a specific module, allowing
// the next notification to pass through immediately.
func (t *Throttle) Reset(module string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSent, module)
}

// ResetAll clears all recorded send times.
func (t *Throttle) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastSent = make(map[string]time.Time)
}
