package diff

import (
	"testing"
)

func makeEncodingEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"production": {
			"API_KEY": "clean-value",
			"SECRET":  "another-clean",
		},
		"staging": {
			"API_KEY": "clean-value",
			"SECRET":  "another-clean",
		},
	}
}

func TestDetectEncodingIssues_NoIssues(t *testing.T) {
	envs := makeEncodingEnvs()
	results := DetectEncodingIssues(envs)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestDetectEncodingIssues_NullByte(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"DB_PASS": "pass\x00word"},
	}
	results := DetectEncodingIssues(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Issue != "null byte detected" {
		t.Errorf("unexpected issue: %s", results[0].Issue)
	}
	if results[0].Key != "DB_PASS" {
		t.Errorf("unexpected key: %s", results[0].Key)
	}
}

func TestDetectEncodingIssues_BOM(t *testing.T) {
	envs := map[string]map[string]string{
		"staging": {"TOKEN": "\xEF\xBB\xBFmy-token"},
	}
	results := DetectEncodingIssues(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Issue != "BOM marker detected" {
		t.Errorf("unexpected issue: %s", results[0].Issue)
	}
}

func TestDetectEncodingIssues_InvalidUTF8(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"API_SECRET": "val\x80ue"},
	}
	results := DetectEncodingIssues(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Issue != "invalid UTF-8 encoding" {
		t.Errorf("unexpected issue: %s", results[0].Issue)
	}
}

func TestDetectEncodingIssues_EmptyEnvs(t *testing.T) {
	results := DetectEncodingIssues(map[string]map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty input, got %d", len(results))
	}
}

func TestDetectEncodingIssues_MultipleEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"KEY": "ok", "BAD": "\x01control"},
		"prod": {"KEY": "ok", "BAD": "clean"},
	}
	results := DetectEncodingIssues(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].EnvName != "dev" {
		t.Errorf("expected env 'dev', got '%s'", results[0].EnvName)
	}
}
