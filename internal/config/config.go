package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full depwatch daemon configuration.
type Config struct {
	PollInterval time.Duration `yaml:"poll_interval"`
	ModulePath   string        `yaml:"module_path"`
	Webhooks     []Webhook     `yaml:"webhooks"`
}

// Webhook defines a single webhook endpoint and its settings.
type Webhook struct {
	Name   string            `yaml:"name"`
	URL    string            `yaml:"url"`
	Secret string            `yaml:"secret"`
	Headers map[string]string `yaml:"headers"`
}

// Load reads and parses the YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.ModulePath == "" {
		return fmt.Errorf("module_path must not be empty")
	}
	if c.PollInterval <= 0 {
		c.PollInterval = 10 * time.Minute
	}
	for i, wh := range c.Webhooks {
		if wh.URL == "" {
			return fmt.Errorf("webhook[%d] url must not be empty", i)
		}
	}
	return nil
}
