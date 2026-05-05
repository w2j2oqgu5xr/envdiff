package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Config holds the parsed CLI configuration for an envdiff run.
type Config struct {
	// Envs maps environment name -> file path, e.g. {"prod": ".env.prod"}
	Envs map[string]string
	// Format is the output format: "text", "csv", or "json".
	Format string
	// ShowMatches controls whether matching keys are printed.
	ShowMatches bool
	// IgnoreKeys is a set of key names to exclude from comparison.
	IgnoreKeys map[string]struct{}
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Envs:        make(map[string]string),
		Format:      "text",
		ShowMatches: false,
		IgnoreKeys:  make(map[string]struct{}),
	}
}

// Validate checks that the config is usable and returns a descriptive error
// if anything is missing or invalid.
func (c *Config) Validate() error {
	if len(c.Envs) < 2 {
		return errors.New("at least two environments must be specified")
	}

	validFormats := map[string]bool{"text": true, "csv": true, "json": true}
	if !validFormats[c.Format] {
		return fmt.Errorf("unknown format %q: must be one of text, csv, json", c.Format)
	}

	for name, path := range c.Envs {
		if strings.TrimSpace(name) == "" {
			return errors.New("environment name must not be empty")
		}
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("file for environment %q not found: %s", name, path)
		}
	}
	return nil
}

// AddIgnoreKeys parses a comma-separated list of key names and adds them to
// the IgnoreKeys set.
func (c *Config) AddIgnoreKeys(raw string) {
	for _, k := range strings.Split(raw, ",") {
		k = strings.TrimSpace(k)
		if k != "" {
			c.IgnoreKeys[k] = struct{}{}
		}
	}
}
