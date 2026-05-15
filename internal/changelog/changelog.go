// Package changelog records and retrieves version change history for monitored modules.
package changelog

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single version change event for a module.
type Entry struct {
	Module  string    `json:"module"`
	OldVer  string    `json:"old_version"`
	NewVer  string    `json:"new_version"`
	DetectedAt time.Time `json:"detected_at"`
}

// Log maintains an ordered list of version change entries backed by a JSON file.
type Log struct {
	mu      sync.RWMutex
	path    string
	entries []Entry
}

// New loads (or creates) a changelog at the given file path.
func New(path string) (*Log, error) {
	l := &Log{path: path}
	if err := l.load(); err != nil {
		return nil, err
	}
	return l, nil
}

// Record appends a new change entry and persists it to disk.
func (l *Log) Record(module, oldVer, newVer string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, Entry{
		Module:     module,
		OldVer:     oldVer,
		NewVer:     newVer,
		DetectedAt: time.Now().UTC(),
	})
	return l.save()
}

// Entries returns a copy of all recorded change entries.
func (l *Log) Entries() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// EntriesFor returns all entries for the given module.
func (l *Log) EntriesFor(module string) []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out []Entry
	for _, e := range l.entries {
		if e.Module == module {
			out = append(out, e)
		}
	}
	return out
}

func (l *Log) load() error {
	data, err := os.ReadFile(l.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &l.entries)
}

func (l *Log) save() error {
	data, err := json.MarshalIndent(l.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o644)
}
