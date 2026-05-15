package checker

import (
	"context"
	"time"

	"github.com/depwatch/internal/cache"
)

// VersionGetter is the interface satisfied by Client and CircuitClient.
type VersionGetter interface {
	LatestVersion(ctx context.Context, module string) (string, error)
}

// CachedClient wraps a VersionGetter with an in-memory TTL cache so that
// repeated polls for the same module within the TTL window skip the network.
type CachedClient struct {
	inner VersionGetter
	cache *cache.Cache
}

// NewCachedClient wraps inner with a cache whose entries live for ttl.
func NewCachedClient(inner VersionGetter, ttl time.Duration) *CachedClient {
	return &CachedClient{
		inner: inner,
		cache: cache.New(ttl),
	}
}

// LatestVersion returns the cached version for module if available and
// unexpired, otherwise delegates to the underlying client and caches the
// result.
func (c *CachedClient) LatestVersion(ctx context.Context, module string) (string, error) {
	if v, ok := c.cache.Get(module); ok {
		return v, nil
	}
	v, err := c.inner.LatestVersion(ctx, module)
	if err != nil {
		return "", err
	}
	c.cache.Set(module, v)
	return v, nil
}

// Invalidate removes the cached entry for module, forcing the next call to
// LatestVersion to fetch a fresh value from upstream.
func (c *CachedClient) Invalidate(module string) {
	c.cache.Delete(module)
}
