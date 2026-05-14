package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultPollInterval = 5 * time.Minute

// Module represents a single Go module to watch.
type Module struct {
	Path string `yaml:"path"`
}

// Config holds the full depwatch configuration.
type Config struct {
	ModulePath   string        `yaml:"module_path"`
	WebhookURL   string        `yaml:"webhook_url"`
	PollInterval time.Duration `yaml:"poll_interval"`
	Modules      []Module      `yaml:"modules"`
}

// ModulePaths returns a deduplicated slice of all module paths in the config.
func (c *Config) ModulePaths() []string {
	seen := make(map[string]struct{}, len(c.Modules))
	paths := make([]string, 0, len(c.Modules))
	for _, m := range c.Modules {
		if _, ok := seen[m.Path]; !ok {
			seen[m.Path] = struct{}{}
			paths = append(paths, m.Path)
		}
	}
	return paths
}

// Load reads and validates a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = defaultPollInterval
	}

	if cfg.WebhookURL == "" {
		return nil, errors.New("config: webhook_url is required")
	}

	if len(cfg.Modules) == 0 && cfg.ModulePath == "" {
		return nil, errors.New("config: at least one module must be specified")
	}

	// Back-fill Modules from legacy module_path field.
	if len(cfg.Modules) == 0 && cfg.ModulePath != "" {
		cfg.Modules = []Module{{Path: cfg.ModulePath}}
	}

	return &cfg, nil
}
