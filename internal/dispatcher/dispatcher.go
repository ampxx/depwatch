// Package dispatcher routes alerts to one or more notifier targets,
// applying per-target filtering, throttling, and severity gating.
package dispatcher

import (
	"context"
	"fmt"
	"log"

	"github.com/yourorg/depwatch/internal/alert"
	"github.com/yourorg/depwatch/internal/filter"
	"github.com/yourorg/depwatch/internal/throttle"
)

// Sender is the interface satisfied by notifier.Notifier.
type Sender interface {
	Notify(ctx context.Context, a alert.Alert) error
}

// Target pairs a Sender with optional filtering and throttling.
type Target struct {
	Name     string
	Sender   Sender
	Filter   *filter.Filter     // nil means allow all
	Throttle *throttle.Throttle // nil means no throttling
	MinSeverity alert.Severity
}

// Dispatcher fans out alerts to all registered targets.
type Dispatcher struct {
	targets []*Target
}

// New returns an empty Dispatcher. Register targets with Add.
func New() *Dispatcher {
	return &Dispatcher{}
}

// Add registers a target with the dispatcher.
func (d *Dispatcher) Add(t *Target) {
	d.targets = append(d.targets, t)
}

// Dispatch sends the alert to every target that accepts it.
// Errors from individual targets are logged but do not abort delivery to others.
func (d *Dispatcher) Dispatch(ctx context.Context, a alert.Alert) {
	for _, t := range d.targets {
		if a.Severity < t.MinSeverity {
			continue
		}
		if t.Filter != nil && !t.Filter.Allow(a.Module) {
			continue
		}
		if t.Throttle != nil && !t.Throttle.Allow(a.Module) {
			continue
		}
		if err := t.Sender.Notify(ctx, a); err != nil {
			log.Printf("dispatcher: target %q failed for module %s: %v", t.Name, a.Module, err)
		}
	}
}

// TargetCount returns the number of registered targets.
func (d *Dispatcher) TargetCount() int {
	return len(d.targets)
}

// Validate checks that all targets have a non-nil Sender.
func (d *Dispatcher) Validate() error {
	for _, t := range d.targets {
		if t.Sender == nil {
			return fmt.Errorf("dispatcher: target %q has nil Sender", t.Name)
		}
	}
	return nil
}
