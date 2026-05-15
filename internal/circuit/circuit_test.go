package circuit

import (
	"testing"
	"time"
)

func TestBreaker_InitiallyClosed(t *testing.T) {
	b := New(3, 50*time.Millisecond)
	if b.State() != StateClosed {
		t.Fatalf("expected StateClosed, got %v", b.State())
	}
	if !b.Allow() {
		t.Fatal("expected Allow() = true on closed breaker")
	}
}

func TestBreaker_OpensAfterThreshold(t *testing.T) {
	b := New(3, 50*time.Millisecond)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != StateClosed {
		t.Fatal("should still be closed after 2 failures")
	}
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatalf("expected StateOpen after threshold, got %v", b.State())
	}
	if b.Allow() {
		t.Fatal("expected Allow() = false when open")
	}
}

func TestBreaker_HalfOpenAfterTimeout(t *testing.T) {
	b := New(1, 30*time.Millisecond)
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatal("expected StateOpen")
	}
	time.Sleep(40 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("expected Allow() = true after timeout")
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", b.State())
	}
}

func TestBreaker_ClosesOnSuccessFromHalfOpen(t *testing.T) {
	b := New(1, 20*time.Millisecond)
	b.RecordFailure()
	time.Sleep(30 * time.Millisecond)
	b.Allow() // transitions to HalfOpen
	b.RecordSuccess()
	if b.State() != StateClosed {
		t.Fatalf("expected StateClosed after success, got %v", b.State())
	}
	if b.failures != 0 {
		t.Fatal("expected failure counter to reset")
	}
}

func TestBreaker_ReOpensOnFailureFromHalfOpen(t *testing.T) {
	b := New(1, 20*time.Millisecond)
	b.RecordFailure()
	time.Sleep(30 * time.Millisecond)
	b.Allow() // transitions to HalfOpen
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatalf("expected StateOpen after failure in HalfOpen, got %v", b.State())
	}
}

func TestBreaker_SuccessResetFailureCount(t *testing.T) {
	b := New(3, 50*time.Millisecond)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	if b.State() != StateClosed {
		t.Fatal("expected StateClosed after success")
	}
	if b.failures != 0 {
		t.Fatalf("expected 0 failures, got %d", b.failures)
	}
}
