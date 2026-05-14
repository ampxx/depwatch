package healthcheck_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/depwatch/internal/healthcheck"
)

func TestHealthy_Returns200(t *testing.T) {
	h := healthcheck.New()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var s healthcheck.Status
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !s.OK {
		t.Error("expected ok=true")
	}
}

func TestUnhealthy_Returns503(t *testing.T) {
	h := healthcheck.New()
	h.SetHealthy(false)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}

	var s healthcheck.Status
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if s.OK {
		t.Error("expected ok=false")
	}
}

func TestSetHealthy_TogglesBothWays(t *testing.T) {
	h := healthcheck.New()

	h.SetHealthy(false)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 after marking unhealthy, got %d", rec.Code)
	}

	h.SetHealthy(true)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 after marking healthy again, got %d", rec.Code)
	}
}

func TestContentType_IsJSON(t *testing.T) {
	h := healthcheck.New()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
}
