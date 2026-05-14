// Package server wires together the HTTP endpoints exposed by depwatch,
// including the health-check and metrics routes.
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/example/depwatch/internal/healthcheck"
	"github.com/example/depwatch/internal/metrics"
)

const (
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultIdleTimeout  = 60 * time.Second
)

// Server wraps an *http.Server and owns the mux.
type Server struct {
	httpServer *http.Server
}

// New constructs a Server that serves the health and metrics endpoints.
// addr should be in the form "host:port" (e.g. ":8080").
func New(addr string, hc *healthcheck.Handler, m *metrics.Metrics) *Server {
	mux := http.NewServeMux()

	mux.Handle("/healthz", hc)
	mux.Handle("/metrics", metrics.NewHandler(m))

	handler := metrics.Middleware(m, mux)

	s := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}

	return &Server{httpServer: s}
}

// Start begins listening and serving HTTP requests. It blocks until the
// provided context is cancelled, at which point it performs a graceful
// shutdown with a 10-second deadline.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("http server error: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	}
}
