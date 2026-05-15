package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/retry"
)

var errTemp = errors.New("temporary error")

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.DefaultConfig(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesAndSucceeds(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}
	err := retry.Do(context.Background(), cfg, func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after eventual success, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}
	err := retry.Do(context.Background(), cfg, func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, retry.ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
		calls++
		return errTemp
	})
	if err == nil {
		t.Fatal("expected error due to cancelled context")
	}
	// Should not have been called because context was already cancelled.
	if calls != 0 {
		t.Fatalf("expected 0 calls, got %d", calls)
	}
}

func TestDo_CancelledBetweenRetries(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	cfg := retry.Config{MaxAttempts: 5, BaseDelay: 50 * time.Millisecond, MaxDelay: time.Second}
	err := retry.Do(ctx, cfg, func() error {
		calls++
		if calls == 1 {
			cancel()
		}
		return errTemp
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
