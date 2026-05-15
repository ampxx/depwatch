package checker

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/circuit"
)

func newCircuitTestClient(srv *httptest.Server, threshold int, timeout time.Duration) *circuitClient {
	c := NewClient(srv.Client(), srv.URL)
	return NewCircuitClient(c, threshold, timeout)
}

func TestCircuitClient_SuccessKeepsClosed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Version":"v1.2.3"}`))
	}))
	defer srv.Close()

	cc := newCircuitTestClient(srv, 3, 50*time.Millisecond)
	v, err := cc.LatestVersion(context.Background(), "example.com/mod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "v1.2.3" {
		t.Fatalf("expected v1.2.3, got %s", v)
	}
	if cc.State() != circuit.StateClosed {
		t.Fatal("expected circuit to remain closed")
	}
}

func TestCircuitClient_OpensAfterThreshold(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	cc := newCircuitTestClient(srv, 2, 50*time.Millisecond)
	for i := 0; i < 2; i++ {
		_, _ = cc.LatestVersion(context.Background(), "example.com/mod")
	}
	if cc.State() != circuit.StateOpen {
		t.Fatalf("expected StateOpen, got %v", cc.State())
	}
}

func TestCircuitClient_RejectsWhenOpen(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	cc := newCircuitTestClient(srv, 1, 50*time.Millisecond)
	_, _ = cc.LatestVersion(context.Background(), "example.com/mod") // trip

	_, err := cc.LatestVersion(context.Background(), "example.com/mod")
	if err == nil {
		t.Fatal("expected error when circuit is open")
	}
	if !errors.Is(err, circuit.ErrOpen) {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestCircuitClient_RecoverAfterTimeout(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Version":"v2.0.0"}`))
	}))
	defer srv.Close()

	cc := newCircuitTestClient(srv, 1, 30*time.Millisecond)
	_, _ = cc.LatestVersion(context.Background(), "example.com/mod") // trip

	time.Sleep(40 * time.Millisecond)

	v, err := cc.LatestVersion(context.Background(), "example.com/mod")
	if err != nil {
		t.Fatalf("expected recovery, got error: %v", err)
	}
	if v != "v2.0.0" {
		t.Fatalf("expected v2.0.0, got %s", v)
	}
	if cc.State() != circuit.StateClosed {
		t.Fatalf("expected StateClosed after recovery, got %v", cc.State())
	}
}
