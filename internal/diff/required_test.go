package diff

import (
	"testing"
)

func makeRequiredEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"production": {
			"DB_HOST": "prod-db",
			"API_KEY": "secret",
			"PORT":    "8080",
		},
		"staging": {
			"DB_HOST": "stage-db",
			"PORT":    "8080",
		},
		"development": {
			"DB_HOST": "localhost",
			"API_KEY": "dev-key",
		},
	}
}

func TestCheckRequired_AllPresent(t *testing.T) {
	envs := makeRequiredEnvs()
	results := CheckRequired(envs, []string{"DB_HOST"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Missing) != 0 {
		t.Errorf("expected no missing envs, got %v", results[0].Missing)
	}
	if len(results[0].Present) != 3 {
		t.Errorf("expected 3 present envs, got %d", len(results[0].Present))
	}
}

func TestCheckRequired_SomeMissing(t *testing.T) {
	envs := makeRequiredEnvs()
	results := CheckRequired(envs, []string{"API_KEY"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	// staging is missing API_KEY
	if len(results[0].Missing) != 1 {
		t.Errorf("expected 1 missing env, got %v", results[0].Missing)
	}
	if results[0].Missing[0] != "staging" {
		t.Errorf("expected staging to be missing, got %s", results[0].Missing[0])
	}
}

func TestCheckRequired_MultipleKeys(t *testing.T) {
	envs := makeRequiredEnvs()
	results := CheckRequired(envs, []string{"API_KEY", "PORT"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestCheckRequired_EmptyEnvs(t *testing.T) {
	results := CheckRequired(nil, []string{"API_KEY"})
	if results != nil {
		t.Errorf("expected nil results for empty envs")
	}
}

func TestCheckRequired_EmptyKeys(t *testing.T) {
	envs := makeRequiredEnvs()
	results := CheckRequired(envs, nil)
	if results != nil {
		t.Errorf("expected nil results for empty required keys")
	}
}

func TestAnyMissing_True(t *testing.T) {
	envs := makeRequiredEnvs()
	results := CheckRequired(envs, []string{"API_KEY", "PORT"})
	if !AnyMissing(results) {
		t.Error("expected AnyMissing to return true")
	}
}

func TestAnyMissing_False(t *testing.T) {
	envs := makeRequiredEnvs()
	results := CheckRequired(envs, []string{"DB_HOST"})
	if AnyMissing(results) {
		t.Error("expected AnyMissing to return false")
	}
}
