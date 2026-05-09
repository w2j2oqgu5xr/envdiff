package output

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeGroupedResults() []diff.CompareResult {
	return []diff.CompareResult{
		{Key: "APP_NAME", Status: "match", Values: map[string]string{"prod": "app", "dev": "app"}},
		{Key: "DB_HOST", Status: "mismatch", Values: map[string]string{"prod": "db.prod", "dev": "localhost"}},
		{Key: "SECRET", Status: "missing", Values: map[string]string{"prod": "s3cr3t", "dev": ""}},
	}
}

func TestGroupedPrinter_EmptyResults(t *testing.T) {
	out := captureGrouped([]diff.CompareResult{}, false)
	if !strings.Contains(out, "No results") {
		t.Errorf("expected 'No results' message, got: %q", out)
	}
}

func TestGroupedPrinter_ShowsMissingSection(t *testing.T) {
	out := captureGrouped(makeGroupedResults(), false)
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING section, got: %q", out)
	}
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET key, got: %q", out)
	}
}

func TestGroupedPrinter_ShowsMismatchSection(t *testing.T) {
	out := captureGrouped(makeGroupedResults(), false)
	if !strings.Contains(out, "MISMATCH") {
		t.Errorf("expected MISMATCH section, got: %q", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST key, got: %q", out)
	}
}

func TestGroupedPrinter_HidesMatchesByDefault(t *testing.T) {
	out := captureGrouped(makeGroupedResults(), false)
	if strings.Contains(out, "APP_NAME") {
		t.Errorf("expected APP_NAME to be hidden when showMatches=false")
	}
}

func TestGroupedPrinter_ShowsMatchesWhenEnabled(t *testing.T) {
	out := captureGrouped(makeGroupedResults(), true)
	if !strings.Contains(out, "APP_NAME") {
		t.Errorf("expected APP_NAME when showMatches=true, got: %q", out)
	}
	if !strings.Contains(out, "MATCH") {
		t.Errorf("expected MATCH section header, got: %q", out)
	}
}

func TestGroupedPrinter_SummaryLine(t *testing.T) {
	out := captureGrouped(makeGroupedResults(), false)
	if !strings.Contains(out, "Summary:") {
		t.Errorf("expected Summary line, got: %q", out)
	}
	if !strings.Contains(out, "total=3") {
		t.Errorf("expected total=3 in summary, got: %q", out)
	}
}

// TestGroupedPrinter_SummaryCountsAllStatuses verifies that the summary line
// correctly reflects the individual counts for each status category.
func TestGroupedPrinter_SummaryCountsAllStatuses(t *testing.T) {
	out := captureGrouped(makeGroupedResults(), false)
	expected := []string{"match=1", "mismatch=1", "missing=1"}
	for _, want := range expected {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in summary, got: %q", want, out)
		}
	}
}

func captureGrouped(results []diff.CompareResult, showMatches bool) string {
	var sb strings.Builder
	p := NewGroupedPrinter(&sb, showMatches)
	p.Print(results)
	return sb.String()
}
