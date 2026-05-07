package watcher

import (
	"context"
	"log"
	"time"

	"github.com/depwatch/internal/checker"
	"github.com/depwatch/internal/config"
	"github.com/depwatch/internal/notifier"
	"github.com/depwatch/internal/store"
)

// Watcher polls dependencies for version changes and sends alerts.
type Watcher struct {
	cfg      *config.Config
	client   *checker.Client
	notifier *notifier.Notifier
	store    *store.Store
}

// New creates a new Watcher with the provided dependencies.
func New(cfg *config.Config, client *checker.Client, n *notifier.Notifier, s *store.Store) *Watcher {
	return &Watcher{
		cfg:      cfg,
		client:   client,
		notifier: n,
		store:    s,
	}
}

// Run starts the polling loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) {
	ticker := time.NewTicker(w.cfg.PollInterval)
	defer ticker.Stop()

	log.Printf("watcher: starting poll loop every %s", w.cfg.PollInterval)
	w.poll(ctx)

	for {
		select {
		case <-ticker.C:
			w.poll(ctx)
		case <-ctx.Done():
			log.Println("watcher: shutting down")
			return
		}
	}
}

// poll checks all configured modules for version updates.
func (w *Watcher) poll(ctx context.Context) {
	for _, mod := range w.cfg.Modules {
		latest, err := w.client.LatestVersion(ctx, mod.Path)
		if err != nil {
			log.Printf("watcher: failed to fetch version for %s: %v", mod.Path, err)
			continue
		}

		prev, err := w.store.Get(mod.Path)
		if err == nil && prev == latest {
			continue
		}

		if err := w.store.Set(mod.Path, latest); err != nil {
			log.Printf("watcher: failed to store version for %s: %v", mod.Path, err)
		}

		if err == nil {
			// prev existed and differs — notify
			log.Printf("watcher: %s updated %s -> %s", mod.Path, prev, latest)
			if nErr := w.notifier.Notify(ctx, mod.Path, prev, latest); nErr != nil {
				log.Printf("watcher: notification failed for %s: %v", mod.Path, nErr)
			}
		} else {
			log.Printf("watcher: %s first seen at %s", mod.Path, latest)
		}
	}
}
