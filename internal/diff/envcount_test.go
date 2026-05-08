package diff

import (
	"testing"
)

func makeEnvcountEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"production": {
			"DB_HOST": "prod-db",
			"DB_PORT": "5432",
			"API_KEY": "secret",
			"LOG_LEVEL": "warn",
		},
		"staging": {
			"DB_HOST": "staging-db",
			"DB_PORT": "5432",
		},
		"development": {
			"DB_HOST": "localhost",
		},
	}
}

func TestCountEnvKeys_ReturnsSortedByName(t *testing.T) {
	envs := makeEnvcountEnvs()
	results := CountEnvKeys(envs)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].EnvName != "development" {
		t.Errorf("expected first env to be development, got %s", results[0].EnvName)
	}
	if results[1].EnvName != "production" {
		t.Errorf("expected second env to be production, got %s", results[1].EnvName)
	}
	if results[2].EnvName != "staging" {
		t.Errorf("expected third env to be staging, got %s", results[2].EnvName)
	}
}

func TestCountEnvKeys_CorrectCounts(t *testing.T) {
	envs := makeEnvcountEnvs()
	results := CountEnvKeys(envs)

	counts := map[string]int{}
	for _, r := range results {
		counts[r.EnvName] = r.KeyCount
	}

	if counts["production"] != 4 {
		t.Errorf("expected production to have 4 keys, got %d", counts["production"])
	}
	if counts["staging"] != 2 {
		t.Errorf("expected staging to have 2 keys, got %d", counts["staging"])
	}
	if counts["development"] != 1 {
		t.Errorf("expected development to have 1 key, got %d", counts["development"])
	}
}

func TestCountEnvKeys_UniqueKeysSorted(t *testing.T) {
	envs := map[string]map[string]string{
		"env": {"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"},
	}
	results := CountEnvKeys(envs)
	keys := results[0].UniqueKeys
	if keys[0] != "A_KEY" || keys[1] != "M_KEY" || keys[2] != "Z_KEY" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

func TestCountEnvKeys_Empty(t *testing.T) {
	results := CountEnvKeys(map[string]map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestKeyCountDeltas_SortedByAbsDelta(t *testing.T) {
	results := []EnvCountResult{
		{EnvName: "production", KeyCount: 10},
		{EnvName: "staging", KeyCount: 7},
		{EnvName: "development", KeyCount: 3},
	}
	deltas := KeyCountDeltas(results)
	if len(deltas) != 3 {
		t.Fatalf("expected 3 deltas, got %d", len(deltas))
	}
	// largest abs delta should be first: production vs development = 7
	if deltas[0].AbsDelta != 7 {
		t.Errorf("expected first delta to be 7, got %d", deltas[0].AbsDelta)
	}
}

func TestKeyCountDeltas_EmptyInput(t *testing.T) {
	deltas := KeyCountDeltas([]EnvCountResult{})
	if len(deltas) != 0 {
		t.Errorf("expected no deltas for empty input")
	}
}
