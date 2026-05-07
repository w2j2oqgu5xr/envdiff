package diff

import (
	"testing"
)

func makeSemverEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"production": {
			"APP_VERSION": "v1.4.0",
			"DB_VERSION":  "5.7.38",
			"NAME":        "myapp",
		},
		"staging": {
			"APP_VERSION": "v1.3.2",
			"DB_VERSION":  "5.7.38",
			"NAME":        "myapp",
		},
		"development": {
			"APP_VERSION": "v1.4.0",
			"DB_VERSION":  "8.0.31",
			"NAME":        "myapp",
		},
	}
}

func TestDetectSemverMismatches_FindsMismatches(t *testing.T) {
	envs := makeSemverEnvs()
	results := DetectSemverMismatches(envs)

	if len(results) != 2 {
		t.Fatalf("expected 2 mismatches, got %d", len(results))
	}

	keys := map[string]bool{}
	for _, r := range results {
		keys[r.Key] = true
	}

	if !keys["APP_VERSION"] {
		t.Error("expected APP_VERSION to be flagged")
	}
	if !keys["DB_VERSION"] {
		t.Error("expected DB_VERSION to be flagged")
	}
}

func TestDetectSemverMismatches_IgnoresNonSemver(t *testing.T) {
	envs := map[string]map[string]string{
		"a": {"NAME": "foo", "APP_VERSION": "v2.0.0"},
		"b": {"NAME": "bar", "APP_VERSION": "v2.0.0"},
	}
	results := DetectSemverMismatches(envs)
	if len(results) != 0 {
		t.Errorf("expected 0 mismatches, got %d", len(results))
	}
}

func TestDetectSemverMismatches_EmptyEnvs(t *testing.T) {
	results := DetectSemverMismatches(map[string]map[string]string{})
	if results != nil {
		t.Errorf("expected nil for empty input, got %v", results)
	}
}

func TestDetectSemverMismatches_SingleEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"APP_VERSION": "v1.0.0"},
	}
	results := DetectSemverMismatches(envs)
	if len(results) != 0 {
		t.Errorf("expected 0 results for single env, got %d", len(results))
	}
}

func TestFormatSemverKey(t *testing.T) {
	m := SemverMismatch{Key: "APP_VERSION", Values: map[string]string{"a": "v1.0.0"}}
	got := FormatSemverKey(m)
	want := "[semver] APP_VERSION"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
