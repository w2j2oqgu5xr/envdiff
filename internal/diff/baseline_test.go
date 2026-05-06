package diff

import (
	"testing"
)

func makeBaselineEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"prod": {
			"DB_HOST": "prod-db",
			"DB_PORT": "5432",
			"SECRET":  "s3cr3t",
		},
		"staging": {
			"DB_HOST": "staging-db",
			"DB_PORT": "5432",
			"DEBUG":   "true",
		},
		"dev": {
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
			"SECRET":  "dev-secret",
			"DEBUG":   "true",
		},
	}
}

func TestCompareToBaseline_InvalidBaseline(t *testing.T) {
	envs := makeBaselineEnvs()
	_, err := CompareToBaseline(envs, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing baseline, got nil")
	}
}

func TestCompareToBaseline_MissingKeys(t *testing.T) {
	envs := makeBaselineEnvs()
	results, err := CompareToBaseline(envs, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var stagingDiff *BaselineDiff
	for i := range results {
		if results[i].Env == "staging" {
			stagingDiff = &results[i]
		}
	}
	if stagingDiff == nil {
		t.Fatal("expected diff for staging")
	}
	if len(stagingDiff.Missing) != 1 || stagingDiff.Missing[0] != "SECRET" {
		t.Errorf("expected SECRET missing, got %v", stagingDiff.Missing)
	}
}

func TestCompareToBaseline_ExtraKeys(t *testing.T) {
	envs := makeBaselineEnvs()
	results, err := CompareToBaseline(envs, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var stagingDiff *BaselineDiff
	for i := range results {
		if results[i].Env == "staging" {
			stagingDiff = &results[i]
		}
	}
	if stagingDiff == nil {
		t.Fatal("expected diff for staging")
	}
	if len(stagingDiff.Extra) != 1 || stagingDiff.Extra[0] != "DEBUG" {
		t.Errorf("expected DEBUG extra, got %v", stagingDiff.Extra)
	}
}

func TestCompareToBaseline_ChangedKeys(t *testing.T) {
	envs := makeBaselineEnvs()
	results, err := CompareToBaseline(envs, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var devDiff *BaselineDiff
	for i := range results {
		if results[i].Env == "dev" {
			devDiff = &results[i]
		}
	}
	if devDiff == nil {
		t.Fatal("expected diff for dev")
	}
	if len(devDiff.Changed) != 1 || devDiff.Changed[0] != "SECRET" {
		t.Errorf("expected SECRET changed, got %v", devDiff.Changed)
	}
}

func TestCompareToBaseline_ResultsAreSorted(t *testing.T) {
	envs := makeBaselineEnvs()
	results, err := CompareToBaseline(envs, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Env != "dev" || results[1].Env != "staging" {
		t.Errorf("expected sorted order dev, staging; got %s, %s", results[0].Env, results[1].Env)
	}
}
