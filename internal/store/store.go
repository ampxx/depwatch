package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Store persists the last known versions of monitored modules.
type Store struct {
	mu       sync.RWMutex
	versions map[string]string
	path     string
}

// New loads or initialises a Store backed by the given file path.
func New(path string) (*Store, error) {
	s := &Store{
		versions: make(map[string]string),
		path:     path,
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading store file: %w", err)
	}

	if err := json.Unmarshal(data, &s.versions); err != nil {
		return nil, fmt.Errorf("parsing store file: %w", err)
	}
	return s, nil
}

// Get returns the stored version for a module, and whether it exists.
func (s *Store) Get(module string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.versions[module]
	return v, ok
}

// Set updates the stored version for a module and persists to disk.
func (s *Store) Set(module, version string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.versions[module] = version
	return s.flush()
}

func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.versions, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling store: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("writing store file: %w", err)
	}
	return nil
}
