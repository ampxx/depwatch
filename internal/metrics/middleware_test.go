package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func errorHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func TestMiddleware_IncrementsChecksTotal(t *testing.T) {
	m := New()
	handler := Middleware(m, http.HandlerFunc(okHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	snap := m.Snapshot()
	if snap.ChecksTotal != 1 {
		t.Errorf("expected ChecksTotal=1, got %d", snap.ChecksTotal)
	}
}

func TestMiddleware_NonOKIncrementsAlertsTotal(t *testing.T) {
	m := New()
	handler := Middleware(m, http.HandlerFunc(errorHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	snap := m.Snapshot()
	if snap.AlertsTotal != 1 {
		t.Errorf("expected AlertsTotal=1, got %d", snap.AlertsTotal)
	}
}

func TestMiddleware_OKDoesNotIncrementAlertsTotal(t *testing.T) {
	m := New()
	handler := Middleware(m, http.HandlerFunc(okHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	snap := m.Snapshot()
	if snap.AlertsTotal != 0 {
		t.Errorf("expected AlertsTotal=0, got %d", snap.AlertsTotal)
	}
}

func TestMiddleware_RecordsLastCheckDuration(t *testing.T) {
	m := New()
	handler := Middleware(m, http.HandlerFunc(okHandler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	snap := m.Snapshot()
	if snap.LastCheckDuration == 0 {
		t.Error("expected LastCheckDuration to be non-zero after request")
	}
}

func TestMiddleware_MultipleRequests(t *testing.T) {
	m := New()
	handler := Middleware(m, http.HandlerFunc(okHandler))

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	snap := m.Snapshot()
	if snap.ChecksTotal != 5 {
		t.Errorf("expected ChecksTotal=5, got %d", snap.ChecksTotal)
	}
}
