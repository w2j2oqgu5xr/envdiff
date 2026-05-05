package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/envdiff/internal/config"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestValidate_TooFewEnvs(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Envs["dev"] = "/some/path"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for fewer than 2 envs")
	}
}

func TestValidate_UnknownFormat(t *testing.T) {
	p1 := writeTempFile(t, ".env.a", "KEY=1\n")
	p2 := writeTempFile(t, ".env.b", "KEY=2\n")
	cfg := config.DefaultConfig()
	cfg.Envs["a"] = p1
	cfg.Envs["b"] = p2
	cfg.Format = "xml"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestValidate_MissingFile(t *testing.T) {
	p1 := writeTempFile(t, ".env.a", "KEY=1\n")
	cfg := config.DefaultConfig()
	cfg.Envs["a"] = p1
	cfg.Envs["b"] = "/nonexistent/.env.b"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_Valid(t *testing.T) {
	p1 := writeTempFile(t, ".env.a", "KEY=1\n")
	p2 := writeTempFile(t, ".env.b", "KEY=2\n")
	cfg := config.DefaultConfig()
	cfg.Envs["a"] = p1
	cfg.Envs["b"] = p2
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddIgnoreKeys(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AddIgnoreKeys("SECRET, TOKEN,DEBUG")
	for _, k := range []string{"SECRET", "TOKEN", "DEBUG"} {
		if _, ok := cfg.IgnoreKeys[k]; !ok {
			t.Errorf("expected key %q in IgnoreKeys", k)
		}
	}
	if len(cfg.IgnoreKeys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(cfg.IgnoreKeys))
	}
}

func TestAddIgnoreKeys_Empty(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AddIgnoreKeys("")
	if len(cfg.IgnoreKeys) != 0 {
		t.Errorf("expected no keys, got %d", len(cfg.IgnoreKeys))
	}
}
