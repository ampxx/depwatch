package checker_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/checker"
)

func newTestClient(server *httptest.Server) *checker.Client {
	return &checker.Client{
		HTTPClient: server.Client(),
		ProxyURL:   server.URL,
	}
}

func TestLatestVersion_Success(t *testing.T) {
	expected := checker.ModuleVersion{
		Version: "v1.2.3",
		Time:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := newTestClient(server)
	got, err := client.LatestVersion("github.com/some/module")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Version != expected.Version {
		t.Errorf("expected version %q, got %q", expected.Version, got.Version)
	}
}

func TestLatestVersion_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.LatestVersion("github.com/missing/module")
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestLatestVersion_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.LatestVersion("github.com/bad/module")
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
