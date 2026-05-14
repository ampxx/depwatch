package server_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/example/depwatch/internal/healthcheck"
	"github.com/example/depwatch/internal/metrics"
	"github.com/example/depwatch/internal/server"
)

func freeAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not find free port: %v", err)
	}
	addr := l.Addr().String()
	_ = l.Close()
	return addr
}

func startServer(t *testing.T) (addr string, cancel context.CancelFunc) {
	t.Helper()

	hc := healthcheck.New()
	hc.SetHealthy(true)
	m := metrics.New()

	addr = freeAddr(t)
	srv := server.New(addr, hc, m)

	ctx, cancel := context.WithCancel(context.Background())

	go func() { _ = srv.Start(ctx) }()

	// Wait until the server is reachable.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	return addr, cancel
}

func TestServer_HealthzReturns200(t *testing.T) {
	addr, cancel := startServer(t)
	defer cancel()

	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestServer_MetricsReturns200(t *testing.T) {
	addr, cancel := startServer(t)
	defer cancel()

	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", addr))
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestServer_GracefulShutdown(t *testing.T) {
	addr, cancel := startServer(t)

	// Cancel triggers shutdown; subsequent requests should fail.
	cancel()
	time.Sleep(200 * time.Millisecond)

	_, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	if err == nil {
		t.Error("expected connection error after shutdown, got none")
	}
}
