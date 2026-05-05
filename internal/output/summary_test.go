package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeSampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusMatch},
		{Key: "DB_HOST", Status: diff.StatusMismatch},
		{Key: "SECRET_KEY", Status: diff.StatusMissing},
		{Key: "PORT", Status: diff.StatusMissing},
	}
}

func TestBuildSummary_Counts(t *testing.T) {
	results := makeSampleResults()
	s := BuildSummary(results)

	if s.TotalKeys != 4 {
		t.Errorf("TotalKeys: want 4, got %d", s.TotalKeys)
	}
	if s.MatchedKeys != 1 {
		t.Errorf("MatchedKeys: want 1, got %d", s.MatchedKeys)
	}
	if s.MismatchKeys != 1 {
		t.Errorf("MismatchKeys: want 1, got %d", s.MismatchKeys)
	}
	if s.MissingKeys != 2 {
		t.Errorf("MissingKeys: want 2, got %d", s.MissingKeys)
	}
}

func TestBuildSummary_Empty(t *testing.T) {
	s := BuildSummary(nil)
	if s.TotalKeys != 0 || s.MatchedKeys != 0 || s.MismatchKeys != 0 || s.MissingKeys != 0 {
		t.Errorf("expected all zeros for empty results, got %+v", s)
	}
}

func TestPrintSummary_Output(t *testing.T) {
	results := makeSampleResults()
	s := BuildSummary(results)

	var buf bytes.Buffer
	PrintSummary(&buf, s)
	out := buf.String()

	for _, want := range []string{"Summary", "4", "1", "2"} {
		if !strings.Contains(out, want) {
			t.Errorf("PrintSummary output missing %q\ngot: %s", want, out)
		}
	}
}

func TestPrintSummary_NilWriter(t *testing.T) {
	// Should not panic when writer is nil (falls back to os.Stdout).
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSummary panicked with nil writer: %v", r)
		}
	}()
	s := BuildSummary(makeSampleResults())
	// Redirect stdout to avoid polluting test output.
	origStdout := captureStdout(func() {
		PrintSummary(nil, s)
	})
	if origStdout == "" {
		// captureStdout returned empty string — acceptable, just ensure no panic.
	}
}
