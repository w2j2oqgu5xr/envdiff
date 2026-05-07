package diff

import (
	"testing"
)

func makePlaceholderEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"production": {
			"API_KEY":    "real-secret-value",
			"DB_HOST":    "prod.db.example.com",
			"AUTH_TOKEN": "tok_abc123",
		},
		"staging": {
			"API_KEY":    "CHANGEME",
			"DB_HOST":    "staging.db.example.com",
			"AUTH_TOKEN": "<your-token-here>",
		},
		"development": {
			"API_KEY":    "TODO",
			"DB_HOST":    "localhost",
			"AUTH_TOKEN": "xxxx1234",
		},
	}
}

func TestDetectPlaceholders_NoPlaceholders(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"KEY": "real-value", "OTHER": "also-real"},
	}
	results := DetectPlaceholders(envs)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestDetectPlaceholders_DetectsChangeme(t *testing.T) {
	envs := map[string]map[string]string{
		"staging": {"API_KEY": "CHANGEME"},
	}
	results := DetectPlaceholders(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "API_KEY" {
		t.Errorf("expected key API_KEY, got %s", results[0].Key)
	}
	if results[0].Pattern != "changeme" {
		t.Errorf("expected pattern changeme, got %s", results[0].Pattern)
	}
}

func TestDetectPlaceholders_DetectsAngleBracket(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"TOKEN": "<your-token>"},
	}
	results := DetectPlaceholders(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].Value != "<your-token>" {
		t.Errorf("unexpected value: %s", results[0].Value)
	}
}

func TestDetectPlaceholders_MultipleEnvs(t *testing.T) {
	envs := makePlaceholderEnvs()
	results := DetectPlaceholders(envs)
	// staging: API_KEY=CHANGEME, AUTH_TOKEN=<your-token-here>
	// development: API_KEY=TODO, AUTH_TOKEN=xxxx1234
	// production: DB_HOST contains "example" — also a match
	if len(results) == 0 {
		t.Fatal("expected placeholder results, got none")
	}
	for _, r := range results {
		if r.Key == "" || r.EnvName == "" || r.Pattern == "" {
			t.Errorf("incomplete result: %+v", r)
		}
	}
}

func TestDetectPlaceholders_EmptyEnvs(t *testing.T) {
	results := DetectPlaceholders(map[string]map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty input, got %d", len(results))
	}
}

func TestDetectPlaceholders_SortedOutput(t *testing.T) {
	envs := map[string]map[string]string{
		"z-env": {"ALPHA": "TODO", "BETA": "CHANGEME"},
		"a-env": {"ALPHA": "fixme"},
	}
	results := DetectPlaceholders(envs)
	for i := 1; i < len(results); i++ {
		prev := results[i-1]
		curr := results[i]
		if prev.Key > curr.Key {
			t.Errorf("results not sorted by key: %s > %s", prev.Key, curr.Key)
		}
		if prev.Key == curr.Key && prev.EnvName > curr.EnvName {
			t.Errorf("results not sorted by env within key: %s > %s", prev.EnvName, curr.EnvName)
		}
	}
}
