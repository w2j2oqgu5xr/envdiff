package diff

import (
	"testing"
)

func makeRenameEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"production": {
			"DB_HOST":     "prod-db",
			"DB_PASSWORD": "secret",
			"API_KEY":     "key123",
		},
		"staging": {
			"DB_HOST":      "staging-db",
			"DATABASE_URL": "postgres://staging",
			"SECRET_KEY":   "stagekey",
		},
	}
}

func TestDetectRenames_NoRenames(t *testing.T) {
	envs := makeRenameEnvs()
	renames := RenameMap{"OLD_KEY": "NEW_KEY"}
	results := DetectRenames(envs, renames)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestDetectRenames_DetectsOldKey(t *testing.T) {
	envs := makeRenameEnvs()
	// production has API_KEY but not API_TOKEN
	renames := RenameMap{"API_KEY": "API_TOKEN"}
	results := DetectRenames(envs, renames)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].OldKey != "API_KEY" || results[0].NewKey != "API_TOKEN" {
		t.Errorf("unexpected rename result: %+v", results[0])
	}
	if results[0].EnvName != "production" {
		t.Errorf("expected env 'production', got %q", results[0].EnvName)
	}
}

func TestDetectRenames_MultipleEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		"alpha": {"OLD_VAR": "val1"},
		"beta":  {"OLD_VAR": "val2"},
	}
	renames := RenameMap{"OLD_VAR": "NEW_VAR"}
	results := DetectRenames(envs, renames)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].EnvName != "alpha" || results[1].EnvName != "beta" {
		t.Errorf("results not sorted by env name: %v", results)
	}
}

func TestDetectRenames_EmptyRenameMap(t *testing.T) {
	envs := makeRenameEnvs()
	results := DetectRenames(envs, RenameMap{})
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty rename map, got %d", len(results))
	}
}

func TestDetectRenames_BothKeysPresent(t *testing.T) {
	// If both old and new key exist, it should NOT be flagged.
	envs := map[string]map[string]string{
		"dev": {"OLD_KEY": "val", "NEW_KEY": "val2"},
	}
	renames := RenameMap{"OLD_KEY": "NEW_KEY"}
	results := DetectRenames(envs, renames)
	if len(results) != 0 {
		t.Errorf("expected 0 results when both keys present, got %d", len(results))
	}
}
