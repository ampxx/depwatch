// Package cache provides a lightweight, thread-safe in-memory TTL cache
// used to reduce redundant upstream requests when polling module versions.
//
// Entries are stored with an expiration time derived from a configurable TTL.
// Expired entries are not evicted automatically; callers should invoke Purge
// periodically or rely on Get returning a cache miss for stale entries.
//
// Example usage:
//
//	c := cache.New(5 * time.Minute)
//	c.Set("github.com/foo/bar", "v1.4.2")
//	if v, ok := c.Get("github.com/foo/bar"); ok {
//		fmt.Println("cached version:", v)
//	}
package cache
