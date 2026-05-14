package metrics_test

import (
	"testing"

	"github.com/example/depwatch/internal/metrics"
)

func TestNew_ZeroValues(t *testing.T) {
	m := metrics.New()
	snap := m.Snapshot()

	if snap.ChecksTotal != 0 {
		t.Errorf("expected ChecksTotal 0, got %d", snap.ChecksTotal)
	}
	if snap.ChangesDetected != 0 {
		t.Errorf("expected ChangesDetected 0, got %d", snap.ChangesDetected)
	}
	if snap.NotificationsSent != 0 {
		t.Errorf("expected NotificationsSent 0, got %d", snap.NotificationsSent)
	}
	if snap.NotificationsFailed != 0 {
		t.Errorf("expected NotificationsFailed 0, got %d", snap.NotificationsFailed)
	}
	if snap.StoreErrors != 0 {
		t.Errorf("expected StoreErrors 0, got %d", snap.StoreErrors)
	}
}

func TestMetrics_Increment(t *testing.T) {
	m := metrics.New()

	m.ChecksTotal.Add(3)
	m.ChangesDetected.Add(1)
	m.NotificationsSent.Add(1)
	m.NotificationsFailed.Add(2)
	m.StoreErrors.Add(1)

	snap := m.Snapshot()

	if snap.ChecksTotal != 3 {
		t.Errorf("expected ChecksTotal 3, got %d", snap.ChecksTotal)
	}
	if snap.ChangesDetected != 1 {
		t.Errorf("expected ChangesDetected 1, got %d", snap.ChangesDetected)
	}
	if snap.NotificationsSent != 1 {
		t.Errorf("expected NotificationsSent 1, got %d", snap.NotificationsSent)
	}
	if snap.NotificationsFailed != 2 {
		t.Errorf("expected NotificationsFailed 2, got %d", snap.NotificationsFailed)
	}
	if snap.StoreErrors != 1 {
		t.Errorf("expected StoreErrors 1, got %d", snap.StoreErrors)
	}
}

func TestSnapshot_IsImmutable(t *testing.T) {
	m := metrics.New()
	m.ChecksTotal.Add(5)

	snap1 := m.Snapshot()
	m.ChecksTotal.Add(10)
	snap2 := m.Snapshot()

	if snap1.ChecksTotal != 5 {
		t.Errorf("snap1 should be 5, got %d", snap1.ChecksTotal)
	}
	if snap2.ChecksTotal != 15 {
		t.Errorf("snap2 should be 15, got %d", snap2.ChecksTotal)
	}
}
