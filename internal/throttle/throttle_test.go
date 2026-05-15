package throttle_test

import (
	"sync"
	"testing"
	"time"

	"github.com/depwatch/internal/throttle"
)

const testModule = "github.com/some/module"

func TestAllow_FirstCallPermitted(t *testing.T) {
	th := throttle.New(5 * time.Second)
	if !th.Allow(testModule) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinWindowSuppressed(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow(testModule)
	if th.Allow(testModule) {
		t.Fatal("expected second call within cooldown to be suppressed")
	}
}

func TestAllow_AfterWindowExpires_Permitted(t *testing.T) {
	now := time.Now()
	th := throttle.New(5 * time.Second)

	// Manually advance time via the internal clock by replacing now func.
	// We use a small cooldown and sleep instead for simplicity.
	th2 := &struct{ *throttle.Throttle }{throttle.New(10 * time.Millisecond)}
	th2.Allow(testModule)
	time.Sleep(20 * time.Millisecond)
	if !th2.Allow(testModule) {
		t.Fatal("expected call after cooldown window to be allowed")
	}
	_ = now
	_ = th
}

func TestAllow_DifferentModulesAreIndependent(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow("github.com/mod/a")
	if !th.Allow("github.com/mod/b") {
		t.Fatal("expected different module to be allowed independently")
	}
}

func TestReset_AllowsImmediateResend(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow(testModule)
	th.Reset(testModule)
	if !th.Allow(testModule) {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAll_ClearsAllModules(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow("github.com/mod/a")
	th.Allow("github.com/mod/b")
	th.ResetAll()
	if !th.Allow("github.com/mod/a") || !th.Allow("github.com/mod/b") {
		t.Fatal("expected all modules to be allowed after ResetAll")
	}
}

func TestAllow_ConcurrentSafe(t *testing.T) {
	th := throttle.New(1 * time.Hour)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			th.Allow(testModule)
		}()
	}
	wg.Wait()
}
