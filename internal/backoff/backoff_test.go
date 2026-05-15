package backoff_test

import (
	"testing"
	"time"

	"github.com/depwatch/internal/backoff"
)

func TestDefaultConfig_Values(t *testing.T) {
	cfg := backoff.DefaultConfig()
	if cfg.Base != 200*time.Millisecond {
		t.Errorf("expected Base 200ms, got %v", cfg.Base)
	}
	if cfg.Max != 30*time.Second {
		t.Errorf("expected Max 30s, got %v", cfg.Max)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier 2.0, got %v", cfg.Multiplier)
	}
	if !cfg.Jitter {
		t.Error("expected Jitter true")
	}
}

func TestNext_GrowsWithAttempts(t *testing.T) {
	cfg := backoff.Config{
		Base:       100 * time.Millisecond,
		Max:        10 * time.Second,
		Multiplier: 2.0,
		Jitter:     false,
	}
	s := backoff.New(cfg)
	prev := s.Next(0)
	for attempt := 1; attempt <= 5; attempt++ {
		next := s.Next(attempt)
		if next <= prev {
			t.Errorf("attempt %d: expected duration > %v, got %v", attempt, prev, next)
		}
		prev = next
	}
}

func TestNext_RespectsMax(t *testing.T) {
	cfg := backoff.Config{
		Base:       1 * time.Second,
		Max:        2 * time.Second,
		Multiplier: 10.0,
		Jitter:     false,
	}
	s := backoff.New(cfg)
	for attempt := 0; attempt <= 10; attempt++ {
		d := s.Next(attempt)
		if d > cfg.Max {
			t.Errorf("attempt %d: duration %v exceeds Max %v", attempt, d, cfg.Max)
		}
	}
}

func TestNext_NegativeAttemptTreatedAsZero(t *testing.T) {
	cfg := backoff.Config{
		Base:       100 * time.Millisecond,
		Max:        5 * time.Second,
		Multiplier: 2.0,
		Jitter:     false,
	}
	s := backoff.New(cfg)
	d := s.Next(-3)
	if d != 100*time.Millisecond {
		t.Errorf("expected 100ms for negative attempt, got %v", d)
	}
}

func TestNext_JitterAddsVariance(t *testing.T) {
	cfg := backoff.Config{
		Base:       500 * time.Millisecond,
		Max:        30 * time.Second,
		Multiplier: 1.0,
		Jitter:     true,
	}
	s := backoff.New(cfg)
	seen := make(map[time.Duration]struct{})
	for i := 0; i < 20; i++ {
		seen[s.Next(0)] = struct{}{}
	}
	if len(seen) < 2 {
		t.Error("expected jitter to produce varying durations")
	}
}

func TestReset_DoesNotPanic(t *testing.T) {
	s := backoff.New(backoff.DefaultConfig())
	s.Reset()
	_ = s.Next(0)
}
