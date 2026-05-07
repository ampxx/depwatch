package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestNotifier(url string) *Notifier {
	n := New(url)
	return n
}

func TestNotify_Success(t *testing.T) {
	var received Payload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := newTestNotifier(ts.URL)
	if err := n.Notify("github.com/foo/bar", "v1.0.0", "v1.1.0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Module != "github.com/foo/bar" {
		t.Errorf("expected module github.com/foo/bar, got %s", received.Module)
	}
	if received.OldVersion != "v1.0.0" {
		t.Errorf("expected old version v1.0.0, got %s", received.OldVersion)
	}
	if received.NewVersion != "v1.1.0" {
		t.Errorf("expected new version v1.1.0, got %s", received.NewVersion)
	}
	if received.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestNotify_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := newTestNotifier(ts.URL)
	err := n.Notify("github.com/foo/bar", "v1.0.0", "v1.1.0")
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestNotify_InvalidURL(t *testing.T) {
	n := newTestNotifier("http://127.0.0.1:0/invalid")
	err := n.Notify("github.com/foo/bar", "v1.0.0", "v1.1.0")
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}
