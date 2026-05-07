package watcher_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/depwatch/internal/checker"
	"github.com/depwatch/internal/config"
	"github.com/depwatch/internal/notifier"
	"github.com/depwatch/internal/store"
	"github.com/depwatch/internal/watcher"
)

func buildWatcher(t *testing.T, proxyURL, webhookURL, storePath string) *watcher.Watcher {
	t.Helper()
	cfg := &config.Config{
		WebhookURL:   webhookURL,
		PollInterval: 50 * time.Millisecond,
		Modules:      []config.Module{{Path: "golang.org/x/text"}},
	}
	client := checker.NewClient(proxyURL)
	n := notifier.New(webhookURL)
	s, err := store.New(storePath)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return watcher.New(cfg, client, n, s)
}

func TestWatcher_DetectsNewVersion(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Version":"v0.3.8"}`))
	}))
	defer proxy.Close()

	notified := make(chan struct{}, 1)
	webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		notified <- struct{}{}
	}))
	defer webhook.Close()

	dir := t.TempDir()
	w := buildWatcher(t, proxy.URL, webhook.URL, filepath.Join(dir, "store.json"))

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	go w.Run(ctx)
	<-ctx.Done()
	// First poll stores version; no notification expected on first seen.
	select {
	case <-notified:
		t.Error("unexpected notification on first version seen")
	default:
	}
}

func TestWatcher_NotifiesOnChange(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Version":"v0.3.9"}`))
	}))
	defer proxy.Close()

	notified := make(chan struct{}, 1)
	webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		notified <- struct{}{}
	}))
	defer webhook.Close()

	dir := t.TempDir()
	storePath := filepath.Join(dir, "store.json")

	// Pre-seed store with old version
	s, _ := store.New(storePath)
	_ = s.Set("golang.org/x/text", "v0.3.8")

	w := buildWatcher(t, proxy.URL, webhook.URL, storePath)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go w.Run(ctx)

	select {
	case <-notified:
		// success
	case <-time.After(400 * time.Millisecond):
		t.Error("expected notification but none received")
	}
}

func TestWatcher_StopsOnContextCancel(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Version":"v1.0.0"}`))
	}))
	defer proxy.Close()

	dir := t.TempDir()
	w := buildWatcher(t, proxy.URL, "http://localhost:9999", filepath.Join(dir, "store.json"))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Error("watcher did not stop after context cancel")
	}
	_ = os.RemoveAll(dir)
}
