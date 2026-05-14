package ratelimit_test

import (
	"sync"
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/ratelimit"
)

func TestAllow_ConsumesTokens(t *testing.T) {
	l := ratelimit.New(3, 1)
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected Allow()=false after bucket empty")
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	now := time.Now()
	l := ratelimit.New(1, 2) // 2 tokens/sec
	l.Allow()                // drain

	// Advance internal clock by injecting a custom now function via Reset trick:
	// Instead, we test the real behaviour with a small sleep.
	_ = now
	time.Sleep(600 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("expected token to refill after 600ms at 2 tok/s")
	}
}

func TestReset_RefilsBucket(t *testing.T) {
	l := ratelimit.New(2, 1)
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("bucket should be empty")
	}
	l.Reset()
	if !l.Allow() {
		t.Fatal("expected Allow()=true after Reset")
	}
}

func TestAllow_ConcurrentSafe(t *testing.T) {
	l := ratelimit.New(100, 1000)
	var wg sync.WaitGroup
	allowed := make(chan bool, 200)
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed <- l.Allow()
		}()
	}
	wg.Wait()
	close(allowed)
	count := 0
	for a := range allowed {
		if a {
			count++
		}
	}
	if count > 100 {
		t.Fatalf("expected at most 100 allowed, got %d", count)
	}
}
