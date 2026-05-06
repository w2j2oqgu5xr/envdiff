package diff

import (
	"testing"
)

func makeDuplicateEnvs() map[string][]string {
	return map[string][]string{
		"production": {
			"DB_HOST=prod.db",
			"DB_HOST=prod2.db",
			"API_KEY=abc",
		},
		"staging": {
			"DB_HOST=stg.db",
			"API_KEY=xyz",
			"API_KEY=xyz2",
			"SECRET=s1",
			"SECRET=s2",
		},
		"development": {
			"DB_HOST=localhost",
			"API_KEY=devkey",
		},
	}
}

func TestDetectDuplicates_NoDuplicates(t *testing.T) {
	raw := map[string][]string{
		"dev": {"A=1", "B=2"},
	}
	results := DetectDuplicates(nil, raw)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestDetectDuplicates_SingleEnvDuplicate(t *testing.T) {
	raw := makeDuplicateEnvs()
	results := DetectDuplicates(nil, raw)

	// production and staging should have duplicates; development should not
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestDetectDuplicates_SortedByEnvName(t *testing.T) {
	raw := makeDuplicateEnvs()
	results := DetectDuplicates(nil, raw)

	if results[0].EnvName != "production" {
		t.Errorf("expected first env to be 'production', got %q", results[0].EnvName)
	}
	if results[1].EnvName != "staging" {
		t.Errorf("expected second env to be 'staging', got %q", results[1].EnvName)
	}
}

func TestDetectDuplicates_CorrectCounts(t *testing.T) {
	raw := makeDuplicateEnvs()
	results := DetectDuplicates(nil, raw)

	var stagingResult *DuplicateResult
	for i := range results {
		if results[i].EnvName == "staging" {
			stagingResult = &results[i]
			break
		}
	}
	if stagingResult == nil {
		t.Fatal("expected staging result")
	}
	if len(stagingResult.Duplicates) != 2 {
		t.Errorf("expected 2 duplicate keys in staging, got %d", len(stagingResult.Duplicates))
	}
}

func TestDetectDuplicates_DuplicatesSortedByKey(t *testing.T) {
	raw := makeDuplicateEnvs()
	results := DetectDuplicates(nil, raw)

	var stagingResult *DuplicateResult
	for i := range results {
		if results[i].EnvName == "staging" {
			stagingResult = &results[i]
			break
		}
	}
	if stagingResult == nil {
		t.Fatal("expected staging result")
	}
	if stagingResult.Duplicates[0].Key != "API_KEY" {
		t.Errorf("expected first duplicate key to be 'API_KEY', got %q", stagingResult.Duplicates[0].Key)
	}
	if stagingResult.Duplicates[1].Key != "SECRET" {
		t.Errorf("expected second duplicate key to be 'SECRET', got %q", stagingResult.Duplicates[1].Key)
	}
}

func TestDetectDuplicates_SkipsCommentsAndBlanks(t *testing.T) {
	raw := map[string][]string{
		"env": {"# comment", "", "KEY=val1", "KEY=val2"},
	}
	results := DetectDuplicates(nil, raw)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Duplicates[0].Count != 2 {
		t.Errorf("expected count 2, got %d", results[0].Duplicates[0].Count)
	}
}
