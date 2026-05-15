package checker

import (
	"context"
	"fmt"
	"time"

	"github.com/yourorg/depwatch/internal/circuit"
)

// circuitClient wraps a Client with a per-host circuit breaker so that
// repeated failures to the module proxy do not flood the network.
type circuitClient struct {
	inner   *Client
	breaker *circuit.Breaker
}

// NewCircuitClient wraps the given Client with a circuit breaker that opens
// after threshold consecutive failures and recovers after timeout.
func NewCircuitClient(c *Client, threshold int, timeout time.Duration) *circuitClient {
	return &circuitClient{
		inner:   c,
		breaker: circuit.New(threshold, timeout),
	}
}

// LatestVersion delegates to the inner client, recording success or failure
// on the circuit breaker accordingly.
func (cc *circuitClient) LatestVersion(ctx context.Context, modulePath string) (string, error) {
	if !cc.breaker.Allow() {
		return "", fmt.Errorf("checker: %w for module %s", circuit.ErrOpen, modulePath)
	}

	version, err := cc.inner.LatestVersion(ctx, modulePath)
	if err != nil {
		cc.breaker.RecordFailure()
		return "", err
	}

	cc.breaker.RecordSuccess()
	return version, nil
}

// State exposes the underlying breaker state for observability.
func (cc *circuitClient) State() circuit.State {
	return cc.breaker.State()
}
