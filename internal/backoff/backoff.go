// Package backoff provides exponential backoff strategies for retry loops.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Strategy defines how wait durations are computed between attempts.
type Strategy interface {
	// Next returns the duration to wait before the nth attempt (0-indexed).
	Next(attempt int) time.Duration
	// Reset resets any internal state.
	Reset()
}

// Config holds parameters for an exponential backoff strategy.
type Config struct {
	// Base is the initial wait duration.
	Base time.Duration
	// Max caps the computed duration.
	Max time.Duration
	// Multiplier is applied on each successive attempt.
	Multiplier float64
	// Jitter adds a random fraction of the computed duration when true.
	Jitter bool
}

// DefaultConfig returns a Config suitable for most network operations.
func DefaultConfig() Config {
	return Config{
		Base:       200 * time.Millisecond,
		Max:        30 * time.Second,
		Multiplier: 2.0,
		Jitter:     true,
	}
}

type exponential struct {
	cfg Config
	rng *rand.Rand
}

// New creates a new exponential backoff Strategy from cfg.
func New(cfg Config) Strategy {
	return &exponential{
		cfg: cfg,
		//nolint:gosec — backoff jitter does not require cryptographic randomness
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (e *exponential) Next(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	base := float64(e.cfg.Base)
	d := base * math.Pow(e.cfg.Multiplier, float64(attempt))
	if e.cfg.Jitter {
		d += e.rng.Float64() * d * 0.3
	}
	result := time.Duration(d)
	if result > e.cfg.Max {
		result = e.cfg.Max
	}
	return result
}

func (e *exponential) Reset() {}
