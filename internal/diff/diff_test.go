package diff

import (
	"testing"
)

func TestCompare_AllMatch(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"APP_PORT": "8080", "DB_HOST": "localhost"},
		"prod": {"APP_PORT": "8080", "DB_HOST": "localhost"},
	}
	res := Compare(envs)
	for _, e := range res.Entries {
		if e.Status != StatusMatch {
			t.Errorf("key %q: expected StatusMatch, got %v", e.Key, e.Status)
		}
	}
}

func TestCompare_MissingKey(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"APP_PORT": "8080", "SECRET": "abc"},
		"prod": {"APP_PORT": "8080"},
	}
	res := Compare(envs)
	for _, e := range res.Entries {
		if e.Key == "SECRET" && e.Status != StatusMissing {
			t.Errorf("expected SECRET to be StatusMissing, got %v", e.Status)
		}
		if e.Key == "APP_PORT" && e.Status != StatusMatch {
			t.Errorf("expected APP_PORT to be StatusMatch, got %v", e.Status)
		}
	}
}

func TestCompare_MismatchedValue(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"DB_HOST": "localhost"},
		"prod": {"DB_HOST": "db.prod.example.com"},
	}
	res := Compare(envs)
	if len(res.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(res.Entries))
	}
	if res.Entries[0].Status != StatusMismatch {
		t.Errorf("expected StatusMismatch, got %v", res.Entries[0].Status)
	}
}

func TestCompare_EnvNamesAreSorted(t *testing.T) {
	envs := map[string]map[string]string{
		"staging": {"X": "1"},
		"dev":     {"X": "1"},
		"prod":    {"X": "1"},
	}
	res := Compare(envs)
	expected := []string{"dev", "prod", "staging"}
	for i, name := range res.EnvNames {
		if name != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], name)
		}
	}
}

func TestCompare_EmptyEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {},
		"prod": {},
	}
	res := Compare(envs)
	if len(res.Entries) != 0 {
		t.Errorf("expected 0 entries for empty envs, got %d", len(res.Entries))
	}
}
