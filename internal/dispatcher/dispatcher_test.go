package dispatcher_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/alert"
	"github.com/yourorg/depwatch/internal/dispatcher"
	"github.com/yourorg/depwatch/internal/filter"
	"github.com/yourorg/depwatch/internal/throttle"
)

type mockSender struct {
	calls atomic.Int32
	err   error
}

func (m *mockSender) Notify(_ context.Context, _ alert.Alert) error {
	m.calls.Add(1)
	return m.err
}

func makeAlert(module, from, to string) alert.Alert {
	a, _ := alert.New(module, from, to)
	return a
}

func TestDispatch_SendsToAllTargets(t *testing.T) {
	d := dispatcher.New()
	s1, s2 := &mockSender{}, &mockSender{}
	d.Add(&dispatcher.Target{Name: "a", Sender: s1})
	d.Add(&dispatcher.Target{Name: "b", Sender: s2})

	d.Dispatch(context.Background(), makeAlert("github.com/foo/bar", "v1.0.0", "v1.1.0"))

	if s1.calls.Load() != 1 || s2.calls.Load() != 1 {
		t.Fatalf("expected both targets called once, got %d %d", s1.calls.Load(), s2.calls.Load())
	}
}

func TestDispatch_FilterBlocksTarget(t *testing.T) {
	d := dispatcher.New()
	s := &mockSender{}
	f := filter.New([]string{"github.com/allowed/*"}, nil)
	d.Add(&dispatcher.Target{Name: "filtered", Sender: s, Filter: f})

	d.Dispatch(context.Background(), makeAlert("github.com/other/pkg", "v1.0.0", "v1.1.0"))

	if s.calls.Load() != 0 {
		t.Fatalf("expected 0 calls, got %d", s.calls.Load())
	}
}

func TestDispatch_ThrottleSuppressesRepeat(t *testing.T) {
	d := dispatcher.New()
	s := &mockSender{}
	th := throttle.New(10 * time.Second)
	d.Add(&dispatcher.Target{Name: "throttled", Sender: s, Throttle: th})

	a := makeAlert("github.com/foo/bar", "v1.0.0", "v1.1.0")
	d.Dispatch(context.Background(), a)
	d.Dispatch(context.Background(), a)

	if s.calls.Load() != 1 {
		t.Fatalf("expected 1 call after throttle, got %d", s.calls.Load())
	}
}

func TestDispatch_SenderErrorContinues(t *testing.T) {
	d := dispatcher.New()
	s1 := &mockSender{err: errors.New("boom")}
	s2 := &mockSender{}
	d.Add(&dispatcher.Target{Name: "fail", Sender: s1})
	d.Add(&dispatcher.Target{Name: "ok", Sender: s2})

	d.Dispatch(context.Background(), makeAlert("github.com/foo/bar", "v1.0.0", "v1.1.0"))

	if s2.calls.Load() != 1 {
		t.Fatalf("second sender should still be called, got %d", s2.calls.Load())
	}
}

func TestDispatch_MinSeverityFilters(t *testing.T) {
	d := dispatcher.New()
	s := &mockSender{}
	d.Add(&dispatcher.Target{Name: "high", Sender: s, MinSeverity: alert.SeverityHigh})

	d.Dispatch(context.Background(), makeAlert("github.com/foo/bar", "v1.0.0", "v1.0.1"))

	if s.calls.Load() != 0 {
		t.Fatalf("low-severity alert should be blocked, got %d calls", s.calls.Load())
	}
}

func TestValidate_NilSenderReturnsError(t *testing.T) {
	d := dispatcher.New()
	d.Add(&dispatcher.Target{Name: "bad", Sender: nil})

	if err := d.Validate(); err == nil {
		t.Fatal("expected error for nil Sender")
	}
}
