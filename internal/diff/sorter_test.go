package diff

import (
	"testing"

	"github.com/user/envdiff/internal/config"
)

func makeResultsForSort() []config.Result {
	return []config.Result{
		{Key: "ZEBRA", Status: config.StatusMatch, Values: map[string]string{"prod": "z", "dev": "z"}},
		{Key: "ALPHA", Status: config.StatusMissing, Values: map[string]string{"prod": "a"}},
		{Key: "BETA", Status: config.StatusMismatch, Values: map[string]string{"prod": "1", "dev": "2"}},
		{Key: "GAMMA", Status: config.StatusMissing, Values: map[string]string{"dev": "g"}},
	}
}

func TestSortResults_ByKey(t *testing.T) {
	results := makeResultsForSort()
	sorted := SortResults(results, SortByKey)

	expected := []string{"ALPHA", "BETA", "GAMMA", "ZEBRA"}
	for i, r := range sorted {
		if r.Key != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, r.Key, expected[i])
		}
	}
}

func TestSortResults_ByStatus(t *testing.T) {
	results := makeResultsForSort()
	sorted := SortResults(results, SortByStatus)

	// Missing first, then mismatch, then match
	if sorted[0].Status != config.StatusMissing || sorted[1].Status != config.StatusMissing {
		t.Errorf("expected first two to be Missing, got %v and %v", sorted[0].Status, sorted[1].Status)
	}
	if sorted[2].Status != config.StatusMismatch {
		t.Errorf("expected index 2 to be Mismatch, got %v", sorted[2].Status)
	}
	if sorted[3].Status != config.StatusMatch {
		t.Errorf("expected index 3 to be Match, got %v", sorted[3].Status)
	}
	// Within Missing, sorted by key
	if sorted[0].Key != "ALPHA" || sorted[1].Key != "GAMMA" {
		t.Errorf("missing keys not sorted: got %q, %q", sorted[0].Key, sorted[1].Key)
	}
}

func TestSortResults_ByEnvCount(t *testing.T) {
	results := makeResultsForSort()
	sorted := SortResults(results, SortByEnvCount)

	// Single-env entries first
	if len(sorted[0].Values) > 1 || len(sorted[1].Values) > 1 {
		t.Errorf("expected single-env entries first")
	}
}

func TestSortResults_DefaultIsKey(t *testing.T) {
	results := makeResultsForSort()
	byDefault := SortResults(results, "")
	byKey := SortResults(results, SortByKey)

	for i := range byDefault {
		if byDefault[i].Key != byKey[i].Key {
			t.Errorf("index %d: default %q != key %q", i, byDefault[i].Key, byKey[i].Key)
		}
	}
}

func TestSortResults_DoesNotMutateOriginal(t *testing.T) {
	results := makeResultsForSort()
	originalFirst := results[0].Key
	SortResults(results, SortByKey)
	if results[0].Key != originalFirst {
		t.Errorf("original slice was mutated")
	}
}
