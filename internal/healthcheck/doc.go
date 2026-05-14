// Package healthcheck exposes an HTTP handler suitable for use as a liveness
// or readiness probe in containerised deployments.
//
// Usage:
//
//	h := healthcheck.New()
//	http.Handle("/healthz", h)
//
//	// Mark unhealthy when a critical subsystem fails.
//	h.SetHealthy(false)
//
// The handler responds with a JSON body containing the current health status,
// daemon uptime, and the time the check was evaluated. It returns HTTP 200
// when healthy and HTTP 503 when unhealthy.
package healthcheck
