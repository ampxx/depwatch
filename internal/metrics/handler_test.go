package metrics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/depwatch/internal/metrics"
)

func TestMetricsHandler_ContentType(t *testing.T) {
	m := metrics.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	m.Handler()(rec, req)

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q; want application/json", got)
	}
}

func TestMetricsHandler_StatusOK(t *testing.T) {
	m := metrics.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	m.Handler()(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
	}
}

func TestMetricsHandler_ReflectsIncrements(t *testing.T) {
	m := metrics.New()
	m.IncChecksTotal()
	m.IncChecksTotal()
	m.IncUpdatesDetected()
	m.IncNotificationsSent()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	m.Handler()(rec, req)

	var snap metrics.Snapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if snap.ChecksTotal != 2 {
		t.Errorf("ChecksTotal = %d; want 2", snap.ChecksTotal)
	}
	if snap.UpdatesDetected != 1 {
		t.Errorf("UpdatesDetected = %d; want 1", snap.UpdatesDetected)
	}
	if snap.NotificationsSent != 1 {
		t.Errorf("NotificationsSent = %d; want 1", snap.NotificationsSent)
	}
	if snap.NotificationErrors != 0 {
		t.Errorf("NotificationErrors = %d; want 0", snap.NotificationErrors)
	}
}

func TestMetricsHandler_ZeroOnFreshMetrics(t *testing.T) {
	m := metrics.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	m.Handler()(rec, req)

	var snap metrics.Snapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if snap.ChecksTotal != 0 || snap.UpdatesDetected != 0 ||
		snap.NotificationsSent != 0 || snap.NotificationErrors != 0 {
		t.Errorf("expected all-zero snapshot, got %+v", snap)
	}
}
