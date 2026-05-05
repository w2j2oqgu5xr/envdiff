package diff

import (
	"testing"
)

func makeGroupResults() []CompareResult {
	return []CompareResult{
		{Key: "APP_NAME", Status: "match", Values: map[string]string{"prod": "app", "dev": "app"}},
		{Key: "DB_HOST", Status: "mismatch", Values: map[string]string{"prod": "db.prod", "dev": "localhost"}},
		{Key: "SECRET", Status: "missing", Values: map[string]string{"prod": "s3cr3t"}},
		{Key: "PORT", Status: "mismatch", Values: map[string]string{"prod": "443", "dev": "8080"}},
		{Key: "LOG_LEVEL", Status: "match", Values: map[string]string{"prod": "info", "dev": "info"}},
	}
}

func TestGroupByStatus_Counts(t *testing.T) {
	results := makeGroupResults()
	groups := GroupByStatus(results)

	if len(groups["match"]) != 2 {
		t.Errorf("expected 2 matches, got %d", len(groups["match"]))
	}
	if len(groups["mismatch"]) != 2 {
		t.Errorf("expected 2 mismatches, got %d", len(groups["mismatch"]))
	}
	if len(groups["missing"]) != 1 {
		t.Errorf("expected 1 missing, got %d", len(groups["missing"]))
	}
}

func TestGroupByStatus_EmptyInput(t *testing.T) {
	groups := GroupByStatus([]CompareResult{})
	for status, g := range groups {
		if len(g) != 0 {
			t.Errorf("expected empty group for %q, got %d", status, len(g))
		}
	}
}

func TestGroupByKey_Keys(t *testing.T) {
	results := makeGroupResults()
	groups := GroupByKey(results)

	if len(groups) != 5 {
		t.Errorf("expected 5 keys, got %d", len(groups))
	}
	if _, ok := groups["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key in group")
	}
}

func TestSummarizeGroups(t *testing.T) {
	results := makeGroupResults()
	groups := GroupByStatus(results)
	s := SummarizeGroups(groups)

	if s.Total != 5 {
		t.Errorf("expected total 5, got %d", s.Total)
	}
	if s.Match != 2 {
		t.Errorf("expected 2 matches, got %d", s.Match)
	}
	if s.Missing != 1 {
		t.Errorf("expected 1 missing, got %d", s.Missing)
	}
	if s.Mismatch != 2 {
		t.Errorf("expected 2 mismatches, got %d", s.Mismatch)
	}
}

func TestSortedGroupKeys(t *testing.T) {
	results := makeGroupResults()
	groups := GroupByStatus(results)
	keys := SortedGroupKeys(groups)

	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	expected := []string{"match", "missing", "mismatch"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("expected key[%d]=%q, got %q", i, expected[i], k)
		}
	}
}
