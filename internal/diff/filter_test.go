package diff_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusMatch},
		{Key: "DB_PASSWORD", Status: diff.StatusMissing},
		{Key: "API_KEY", Status: diff.StatusMismatch},
		{Key: "SECRET_TOKEN", Status: diff.StatusMissing},
	}
}

func TestApplyIgnore_NoIgnoreKeys(t *testing.T) {
	results := makeResults()
	got := diff.ApplyIgnore(results, nil)
	if len(got) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(got))
	}
}

func TestApplyIgnore_SingleKey(t *testing.T) {
	results := makeResults()
	got := diff.ApplyIgnore(results, []string{"DB_PASSWORD"})
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
	for _, r := range got {
		if r.Key == "DB_PASSWORD" {
			t.Error("DB_PASSWORD should have been filtered out")
		}
	}
}

func TestApplyIgnore_MultipleKeys(t *testing.T) {
	results := makeResults()
	got := diff.ApplyIgnore(results, []string{"API_KEY", "SECRET_TOKEN"})
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestApplyIgnore_CaseInsensitive(t *testing.T) {
	results := makeResults()
	got := diff.ApplyIgnore(results, []string{"app_name", "api_key"})
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	for _, r := range got {
		if r.Key == "APP_NAME" || r.Key == "API_KEY" {
			t.Errorf("key %q should have been filtered out", r.Key)
		}
	}
}

func TestApplyIgnore_EmptyResults(t *testing.T) {
	got := diff.ApplyIgnore([]diff.Result{}, []string{"APP_NAME"})
	if len(got) != 0 {
		t.Errorf("expected empty results, got %d", len(got))
	}
}

func TestApplyIgnore_AllFiltered(t *testing.T) {
	results := makeResults()
	keys := []string{"APP_NAME", "DB_PASSWORD", "API_KEY", "SECRET_TOKEN"}
	got := diff.ApplyIgnore(results, keys)
	if len(got) != 0 {
		t.Errorf("expected 0 results, got %d", len(got))
	}
}
