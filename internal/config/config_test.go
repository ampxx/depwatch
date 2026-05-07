package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "depwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	raw := `
module_path: /srv/myapp
poll_interval: 5m
webhooks:
  - name: slack
    url: https://hooks.slack.com/xxx
    secret: topsecret
`
	path := writeTempConfig(t, raw)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ModulePath != "/srv/myapp" {
		t.Errorf("module_path = %q, want /srv/myapp", cfg.ModulePath)
	}
	if cfg.PollInterval != 5*time.Minute {
		t.Errorf("poll_interval = %v, want 5m", cfg.PollInterval)
	}
	if len(cfg.Webhooks) != 1 || cfg.Webhooks[0].Name != "slack" {
		t.Errorf("unexpected webhooks: %+v", cfg.Webhooks)
	}
}

func TestLoad_DefaultPollInterval(t *testing.T) {
	raw := `module_path: /srv/app\nwebhooks:\n  - url: https://example.com/hook\n`
	path := writeTempConfig(t, "module_path: /srv/app\nwebhooks:\n  - url: https://example.com/hook\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PollInterval != 10*time.Minute {
		t.Errorf("default poll_interval = %v, want 10m", cfg.PollInterval)
	}
	_ = raw
}

func TestLoad_MissingModulePath(t *testing.T) {
	path := writeTempConfig(t, "webhooks:\n  - url: https://example.com/hook\n")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for missing module_path")
	}
}

func TestLoad_MissingWebhookURL(t *testing.T) {
	path := writeTempConfig(t, "module_path: /srv/app\nwebhooks:\n  - name: broken\n")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for webhook with empty url")
	}
}
