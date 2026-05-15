// Package retry provides a simple exponential backoff retry mechanism
// for use when making outbound HTTP requests to external services.
package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// ErrMaxAttempts is returned when all retry attempts have been exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Config holds the configuration for the retry mechanism.
type Config struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// BaseDelay is the initial delay between retries.
	BaseDelay time.Duration
	// MaxDelay caps the exponential backoff delay.
	MaxDelay time.Duration
}

// DefaultConfig returns a sensible default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    10 * time.Second,
	}
}

// Do executes fn up to cfg.MaxAttempts times, backing off exponentially
// between attempts. It returns nil on the first success, or ErrMaxAttempts
// if all attempts fail. The context is respected between retries.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if attempt == cfg.MaxAttempts-1 {
			break
		}
		delay := delay(cfg.BaseDelay, cfg.MaxDelay, attempt)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	_ = lastErr
	return ErrMaxAttempts
}

// delay computes the exponential backoff for a given attempt, capped at max.
func delay(base, max time.Duration, attempt int) time.Duration {
	d := time.Duration(float64(base) * math.Pow(2, float64(attempt)))
	if d > max {
		return max
	}
	return d
}
