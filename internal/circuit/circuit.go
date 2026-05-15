// Package circuit implements a simple circuit breaker to prevent
// repeated calls to failing downstream services.
package circuit

import (
	"errors"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// ErrOpen is returned when the circuit breaker is open.
var ErrOpen = errors.New("circuit breaker is open")

// Breaker is a circuit breaker that trips after a threshold of consecutive
// failures and recovers after a configurable timeout.
type Breaker struct {
	mu           sync.Mutex
	failures     int
	threshold    int
	timeout      time.Duration
	lastFailure  time.Time
	state        State
}

// New creates a new Breaker with the given failure threshold and recovery timeout.
func New(threshold int, timeout time.Duration) *Breaker {
	return &Breaker{
		threshold: threshold,
		timeout:   timeout,
		state:     StateClosed,
	}
}

// Allow reports whether a call should be allowed through.
// It transitions the breaker from Open to HalfOpen once the timeout elapses.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(b.lastFailure) >= b.timeout {
			b.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess records a successful call, closing the circuit if it was half-open.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure records a failed call, potentially opening the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	b.lastFailure = time.Now()
	if b.failures >= b.threshold {
		b.state = StateOpen
	}
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
