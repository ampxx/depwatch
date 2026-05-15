package dedupe_test

import (
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/dedupe"
)

func TestIsDuplicate_FirstCallNotDuplicate(t *testing.T) {
	d := dedupe.New(time.Minute)
	if d.IsDuplicate("github.com/foo/bar", "v1.2.0") {
		t.Fatal("first call should not be a duplicate")
	}
}

func TestIsDuplicate_SecondCallWithinWindowIsDuplicate(t *testing.T) {
	d := dedupe.New(time.Minute)
	d.IsDuplicate("github.com/foo/bar", "v1.2.0")
	if !d.IsDuplicate("github.com/foo/bar", "v1.2.0") {
		t.Fatal("second call within window should be a duplicate")
	}
}

func TestIsDuplicate_DifferentVersionNotDuplicate(t *testing.T) {
	d := dedupe.New(time.Minute)
	d.IsDuplicate("github.com/foo/bar", "v1.2.0")
	if d.IsDuplicate("github.com/foo/bar", "v1.3.0") {
		t.Fatal("different version should not be a duplicate")
	}
}

func TestIsDuplicate_DifferentModuleNotDuplicate(t *testing.T) {
	d := dedupe.New(time.Minute)
	d.IsDuplicate("github.com/foo/bar", "v1.0.0")
	if d.IsDuplicate("github.com/foo/baz", "v1.0.0") {
		t.Fatal("different module should not be a duplicate")
	}
}

func TestIsDuplicate_AfterWindowExpires_NotDuplicate(t *testing.T) {
	d := dedupe.New(10 * time.Millisecond)
	d.IsDuplicate("github.com/foo/bar", "v1.0.0")
	time.Sleep(20 * time.Millisecond)
	if d.IsDuplicate("github.com/foo/bar", "v1.0.0") {
		t.Fatal("call after window expiry should not be a duplicate")
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	d := dedupe.New(10 * time.Millisecond)
	d.IsDuplicate("github.com/foo/bar", "v1.0.0")
	d.IsDuplicate("github.com/foo/baz", "v2.0.0")

	time.Sleep(20 * time.Millisecond)
	d.Purge()

	if got := d.Len(); got != 0 {
		t.Fatalf("expected 0 entries after purge, got %d", got)
	}
}

func TestLen_CountsActiveEntries(t *testing.T) {
	d := dedupe.New(time.Minute)
	d.IsDuplicate("github.com/a/b", "v1.0.0")
	d.IsDuplicate("github.com/c/d", "v2.0.0")

	if got := d.Len(); got != 2 {
		t.Fatalf("expected 2 active entries, got %d", got)
	}
}
