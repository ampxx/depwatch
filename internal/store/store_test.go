package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/depwatch/internal/store"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "versions.json")
}

func TestStore_SetAndGet(t *testing.T) {
	s, err := store.New(tempStorePath(t))
	if err != nil {
		t.Fatalf("unexpected error creating store: %v", err)
	}

	if err := s.Set("github.com/foo/bar", "v1.0.0"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, ok := s.Get("github.com/foo/bar")
	if !ok {
		t.Fatal("expected module to be present")
	}
	if got != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %q", got)
	}
}

func TestStore_MissingModule(t *testing.T) {
	s, err := store.New(tempStorePath(t))
	if err != nil {
		t.Fatalf("unexpected error creating store: %v", err)
	}

	_, ok := s.Get("github.com/missing/module")
	if ok {
		t.Error("expected module to be absent")
	}
}

func TestStore_PersistsAcrossReload(t *testing.T) {
	path := tempStorePath(t)

	s1, err := store.New(path)
	if err != nil {
		t.Fatalf("creating first store: %v", err)
	}
	if err := s1.Set("github.com/persist/me", "v2.3.4"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	s2, err := store.New(path)
	if err != nil {
		t.Fatalf("reloading store: %v", err)
	}
	got, ok := s2.Get("github.com/persist/me")
	if !ok {
		t.Fatal("expected module after reload")
	}
	if got != "v2.3.4" {
		t.Errorf("expected v2.3.4 after reload, got %q", got)
	}
}

func TestStore_InvalidFile(t *testing.T) {
	path := tempStorePath(t)
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatalf("writing bad file: %v", err)
	}
	_, err := store.New(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON file")
	}
}
