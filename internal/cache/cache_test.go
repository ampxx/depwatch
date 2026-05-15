package cache_test

import (
	"testing"
	"time"

	"github.com/depwatch/internal/cache"
)

func TestGet_MissingKey(t *testing.T) {
	c := cache.New(time.Minute)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected miss for unknown key")
	}
}

func TestSet_AndGet(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("mod", "v1.2.3")
	v, ok := c.Get("mod")
	if !ok {
		t.Fatal("expected hit after Set")
	}
	if v != "v1.2.3" {
		t.Fatalf("expected v1.2.3, got %s", v)
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("mod", "v1.0.0")
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("mod")
	if ok {
		t.Fatal("expected miss for expired entry")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("mod", "v1.0.0")
	c.Delete("mod")
	_, ok := c.Get("mod")
	if ok {
		t.Fatal("expected miss after Delete")
	}
}

func TestPurge_RemovesExpired(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("a", "v1")
	c.Set("b", "v2")
	time.Sleep(20 * time.Millisecond)
	c.Set("c", "v3") // fresh
	c.Purge()
	if c.Len() != 1 {
		t.Fatalf("expected 1 entry after purge, got %d", c.Len())
	}
	_, ok := c.Get("c")
	if !ok {
		t.Fatal("fresh entry should survive purge")
	}
}

func TestLen_ReflectsCount(t *testing.T) {
	c := cache.New(time.Minute)
	if c.Len() != 0 {
		t.Fatal("expected empty cache")
	}
	c.Set("x", "1")
	c.Set("y", "2")
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}
