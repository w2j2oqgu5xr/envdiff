package diff

import (
	"testing"
)

func makeWhitespaceEnvs(data map[string]map[string]string) map[string]map[string]string {
	return data
}

func TestDetectWhitespaceDiffs_NoIssues(t *testing.T) {
	envs := makeWhitespaceEnvs(map[string]map[string]string{
		"prod": {"HOST": "localhost", "PORT": "8080"},
		"dev":  {"HOST": "localhost", "PORT": "8080"},
	})
	results := DetectWhitespaceDiffs(envs)
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestDetectWhitespaceDiffs_DetectsLeadingSpace(t *testing.T) {
	envs := makeWhitespaceEnvs(map[string]map[string]string{
		"prod": {"HOST": " localhost"},
		"dev":  {"HOST": "localhost"},
	})
	results := DetectWhitespaceDiffs(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "HOST" {
		t.Errorf("expected key HOST, got %s", results[0].Key)
	}
	if results[0].Trimmed != "localhost" {
		t.Errorf("expected trimmed 'localhost', got %q", results[0].Trimmed)
	}
}

func TestDetectWhitespaceDiffs_DetectsTrailingSpace(t *testing.T) {
	envs := makeWhitespaceEnvs(map[string]map[string]string{
		"prod": {"API_KEY": "abc123  "},
		"dev":  {"API_KEY": "abc123"},
	})
	results := DetectWhitespaceDiffs(envs)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].RawValue != "abc123  " {
		t.Errorf("unexpected raw value: %q", results[0].RawValue)
	}
}

func TestDetectWhitespaceDiffs_EmptyEnvs(t *testing.T) {
	results := DetectWhitespaceDiffs(map[string]map[string]string{})
	if results != nil {
		t.Errorf("expected nil results for empty input")
	}
}

func TestDetectWhitespaceDiffs_SingleEnvSkipped(t *testing.T) {
	// With only one env, there's nothing to compare against.
	envs := makeWhitespaceEnvs(map[string]map[string]string{
		"prod": {"HOST": " localhost "},
	})
	results := DetectWhitespaceDiffs(envs)
	if len(results) != 0 {
		t.Errorf("expected no results for single env, got %d", len(results))
	}
}

func TestDetectWhitespaceDiffs_MultipleKeys(t *testing.T) {
	envs := makeWhitespaceEnvs(map[string]map[string]string{
		"prod": {"HOST": " localhost", "PORT": "8080 ", "DB": "mydb"},
		"dev":  {"HOST": "localhost", "PORT": "8080", "DB": "mydb"},
	})
	results := DetectWhitespaceDiffs(envs)
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}
