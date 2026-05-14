package diff

import (
	"testing"
)

func makeSchemaEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"production": {
			"DB_HOST": "prod-db",
			"API_KEY": "abc123",
			"LOG_LEVEL": "warn",
		},
		"staging": {
			"DB_HOST":  "staging-db",
			"API_KEY":  "xyz789",
		},
	}
}

func TestValidateSchema_AllPresent(t *testing.T) {
	envs := makeSchemaEnvs()
	results := ValidateSchema(envs, []string{"DB_HOST", "API_KEY"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		for env, present := range r.Present {
			if !present {
				t.Errorf("expected key %q present in %q", r.Key, env)
			}
		}
	}
}

func TestValidateSchema_MissingKey(t *testing.T) {
	envs := makeSchemaEnvs()
	results := ValidateSchema(envs, []string{"LOG_LEVEL"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Key != "LOG_LEVEL" {
		t.Errorf("unexpected key %q", r.Key)
	}
	if !r.Present["production"] {
		t.Error("expected LOG_LEVEL present in production")
	}
	if r.Present["staging"] {
		t.Error("expected LOG_LEVEL absent in staging")
	}
}

func TestValidateSchema_SortedResults(t *testing.T) {
	envs := makeSchemaEnvs()
	results := ValidateSchema(envs, []string{"LOG_LEVEL", "API_KEY", "DB_HOST"})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Key != "API_KEY" || results[1].Key != "DB_HOST" || results[2].Key != "LOG_LEVEL" {
		t.Errorf("results not sorted: %v", []string{results[0].Key, results[1].Key, results[2].Key})
	}
}

func TestValidateSchema_EmptyEnvs(t *testing.T) {
	results := ValidateSchema(nil, []string{"DB_HOST"})
	if results != nil {
		t.Errorf("expected nil results for empty envs")
	}
}

func TestValidateSchema_EmptySchema(t *testing.T) {
	envs := makeSchemaEnvs()
	results := ValidateSchema(envs, nil)
	if results != nil {
		t.Errorf("expected nil results for empty schema")
	}
}

func TestSchemaFullyCovered_True(t *testing.T) {
	envs := makeSchemaEnvs()
	results := ValidateSchema(envs, []string{"DB_HOST", "API_KEY"})
	if !SchemaFullyCovered(results) {
		t.Error("expected schema to be fully covered")
	}
}

func TestSchemaFullyCovered_False(t *testing.T) {
	envs := makeSchemaEnvs()
	results := ValidateSchema(envs, []string{"LOG_LEVEL"})
	if SchemaFullyCovered(results) {
		t.Error("expected schema not fully covered")
	}
}
