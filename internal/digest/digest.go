// Package digest provides utilities for computing and comparing
// content hashes of dependency manifests, enabling depwatch to detect
// changes in go.sum or go.mod files without re-fetching every module.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
)

// Digester computes and caches SHA-256 digests of files on disk.
type Digester struct {
	mu    sync.Mutex
	cache map[string]string // path -> last known hex digest
}

// New returns an initialised Digester.
func New() *Digester {
	return &Digester{
		cache: make(map[string]string),
	}
}

// Changed reports whether the file at path has a different SHA-256 digest
// compared to the last call to Changed (or HashFile) for the same path.
// The first call for a given path always returns false and seeds the cache.
func (d *Digester) Changed(path string) (bool, error) {
	current, err := HashFile(path)
	if err != nil {
		return false, fmt.Errorf("digest: hashing %q: %w", path, err)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	prev, seen := d.cache[path]
	d.cache[path] = current

	if !seen {
		return false, nil
	}
	return prev != current, nil
}

// Reset clears the cached digest for path so the next call to Changed
// treats it as unseen.
func (d *Digester) Reset(path string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.cache, path)
}

// HashFile computes the SHA-256 hex digest of the file at path.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
