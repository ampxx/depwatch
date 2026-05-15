// Package circuit provides a circuit breaker implementation for depwatch.
//
// A circuit breaker prevents repeated calls to a failing downstream service
// (such as a webhook endpoint or the Go module proxy) by "opening" after a
// configurable number of consecutive failures. Once open, all calls are
// rejected immediately with ErrOpen. After a recovery timeout the breaker
// moves to the half-open state and allows a single probe call through. A
// successful probe closes the circuit; a failed probe reopens it.
//
// Usage:
//
//	br := circuit.New(5, 30*time.Second)
//	if !br.Allow() {
//	    return circuit.ErrOpen
//	}
//	if err := doCall(); err != nil {
//	    br.RecordFailure()
//	    return err
//	}
//	br.RecordSuccess()
package circuit
