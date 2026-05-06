package diff

import (
	"testing"
)

func makeUnusedEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"production": {
			"DB_HOST": "prod-db",
			"API_KEY": "abc123",
			"LEGACY_FLAG": "true",
		},
		"staging": {
			"DB_HOST": "staging-db",
			"API_KEY": "xyz789",
		},
		"development": {
			"DB_HOST": "localhost",
			"API_KEY": "devkey",
			"DEBUG": "true",
		},
	}
}

func TestDetectUnused_NoUnused(t *testing.T) {
	envs := map[string]map[string]string{
		"a": {"FOO": "1", "BAR": "2"},
		"b": {"FOO": "1", "BAR": "2"},
	}
	results := DetectUnused(envs)
	if len(results) != 0 {
		t.Errorf("expected no unused keys, got %d", len(results))
	}
}

func TestDetectUnused_SingleEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"only": {"FOO": "bar"},
	}
	results := DetectUnused(envs)
	if len(results) != 0 {
		t.Errorf("expected no results for single env, got %d", len(results))
	}
}

func TestDetectUnused_DetectsPartialKeys(t *testing.T) {
	envs := makeUnusedEnvs()
	results := DetectUnused(envs)

	// LEGACY_FLAG and DEBUG are only in one env each
	if len(results) < 2 {
		t.Fatalf("expected at least 2 unused results, got %d", len(results))
	}

	keyMap := make(map[string]UnusedResult)
	for _, r := range results {
		keyMap[r.Key] = r
	}

	if _, ok := keyMap["LEGACY_FLAG"]; !ok {
		t.Error("expected LEGACY_FLAG to be detected as unused")
	}
	if _, ok := keyMap["DEBUG"]; !ok {
		t.Error("expected DEBUG to be detected as unused")
	}
}

func TestDetectUnused_AbsentFromIsCorrect(t *testing.T) {
	envs := makeUnusedEnvs()
	results := DetectUnused(envs)

	keyMap := make(map[string]UnusedResult)
	for _, r := range results {
		keyMap[r.Key] = r
	}

	legacy := keyMap["LEGACY_FLAG"]
	if len(legacy.AbsentFrom) != 2 {
		t.Errorf("expected LEGACY_FLAG absent from 2 envs, got %d", len(legacy.AbsentFrom))
	}
	if len(legacy.PresentIn) != 1 || legacy.PresentIn[0] != "production" {
		t.Errorf("expected LEGACY_FLAG present only in production, got %v", legacy.PresentIn)
	}
}

func TestDetectUnused_ResultsSortedByKey(t *testing.T) {
	envs := makeUnusedEnvs()
	results := DetectUnused(envs)

	for i := 1; i < len(results); i++ {
		if results[i].Key < results[i-1].Key {
			t.Errorf("results not sorted: %s before %s", results[i-1].Key, results[i].Key)
		}
	}
}
