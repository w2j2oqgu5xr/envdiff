package config

import (
	"errors"
	"fmt"
	"os"
)

// Status represents the comparison state of a key across environments.
type Status string

const (
	StatusMatch    Status = "match"
	StatusMismatch Status = "mismatch"
	StatusMissing  Status = "missing"
)

// Result holds the comparison outcome for a single key.
type Result struct {
	Key    string
	Status Status
	Values map[string]string // env name -> value (empty string if missing)
}

// Config holds all runtime configuration for envdiff.
type Config struct {
	Envs        []EnvFile
	ShowMatches bool
	Format      string
	IgnoreKeys  []string
	SortOrder   string
}

// EnvFile pairs an environment name with its file path.
type EnvFile struct {
	Name string
	Path string
}

var validFormats = map[string]bool{
	"text": true,
	"csv":  true,
	"json": true,
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Format:    "text",
		SortOrder: "key",
	}
}

// Validate checks that the Config is valid before use.
func (c *Config) Validate() error {
	if len(c.Envs) < 2 {
		return errors.New("at least two environments are required")
	}
	if !validFormats[c.Format] {
		return fmt.Errorf("unknown format %q: must be one of text, csv, json", c.Format)
	}
	for _, env := range c.Envs {
		if _, err := os.Stat(env.Path); os.IsNotExist(err) {
			return fmt.Errorf("file not found for env %q: %s", env.Name, env.Path)
		}
	}
	return nil
}
