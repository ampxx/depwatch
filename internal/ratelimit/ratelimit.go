// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how frequently webhook notifications are dispatched.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a thread-safe token-bucket rate limiter.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	now      func() time.Time
}

// New creates a Limiter that allows up to max events with a refill rate of
// rate tokens per second. max must be >= 1 and rate must be > 0.
func New(max int, rate float64) *Limiter {
	return &Limiter{
		tokens:   float64(max),
		max:      float64(max),
		rate:     rate,
		lastTick: time.Now(),
		now:      time.Now,
	}
}

// Allow reports whether an event may proceed. It consumes one token when true.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Reset refills the bucket to its maximum capacity.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tokens = l.max
	l.lastTick = l.now()
}
