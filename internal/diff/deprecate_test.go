package diff

import (
	"testing"
)

func makeDeprecateEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"production": {
			"DATABASE_URL": "postgres://prod",
			"OLD_API_KEY":  "secret",
			"NEW_API_KEY":  "newsecret",
		},
		"staging": {
			"DATABASE_URL": "postgres://staging",
			"OLD_API_KEY":  "stagingsecret",
		},
		"development": {
			"DATABASE_URL": "postgres://dev",
			"NEW_API_KEY":  "devkey",
		},
	}
}

func TestDetectDeprecated_NoMatches(t *testing.T) {
	envs := makeDeprecateEnvs()
	depMap := map[string]string{"LEGACY_TOKEN": "use TOKEN instead"}
	results := DetectDeprecated(envs, depMap)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestDetectDeprecated_SingleKey(t *testing.T) {
	envs := makeDeprecateEnvs()
	depMap := map[string]string{"OLD_API_KEY": "use NEW_API_KEY instead"}
	results := DetectDeprecated(envs, depMap)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "OLD_API_KEY" {
		t.Errorf("unexpected key %q", results[0].Key)
	}
	if results[0].Message != "use NEW_API_KEY instead" {
		t.Errorf("unexpected message %q", results[0].Message)
	}
	if len(results[0].Envs) != 2 {
		t.Errorf("expected 2 envs, got %d", len(results[0].Envs))
	}
}

func TestDetectDeprecated_CaseInsensitiveKey(t *testing.T) {
	envs := makeDeprecateEnvs()
	depMap := map[string]string{"old_api_key": "use NEW_API_KEY instead"}
	results := DetectDeprecated(envs, depMap)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestDetectDeprecated_EmptyEnvs(t *testing.T) {
	results := DetectDeprecated(nil, map[string]string{"OLD": "gone"})
	if results != nil {
		t.Errorf("expected nil results for empty envs")
	}
}

func TestDetectDeprecated_EmptyDeprecationMap(t *testing.T) {
	results := DetectDeprecated(makeDeprecateEnvs(), nil)
	if results != nil {
		t.Errorf("expected nil results for empty deprecation map")
	}
}

func TestDetectDeprecated_SortedOutput(t *testing.T) {
	envs := makeDeprecateEnvs()
	depMap := map[string]string{
		"NEW_API_KEY": "will be removed in v2",
		"OLD_API_KEY": "use NEW_API_KEY instead",
	}
	results := DetectDeprecated(envs, depMap)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Key >= results[1].Key {
		t.Errorf("results not sorted: %q >= %q", results[0].Key, results[1].Key)
	}
}
