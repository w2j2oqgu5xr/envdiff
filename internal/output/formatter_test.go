package output

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{
			Key:    "DB_HOST",
			Status: diff.StatusMismatch,
			Values: map[string]string{"prod": "db.prod.example.com", "staging": "db.staging.example.com"},
		},
		{
			Key:    "SECRET_KEY",
			Status: diff.StatusMissing,
			Values: map[string]string{"prod": "abc123"},
		},
	}
}

func TestFormatResults_UnknownFormat(t *testing.T) {
	_, err := FormatResults(sampleResults(), Format("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestFormatResults_EmptyText(t *testing.T) {
	out, err := FormatResults(nil, FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected 'No differences' message, got: %q", out)
	}
}

func TestFormatResults_Text(t *testing.T) {
	out, err := FormatResults(sampleResults(), FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[MISMATCH]") {
		t.Errorf("expected [MISMATCH] in output, got: %q", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output")
	}
}

func TestFormatResults_CSV(t *testing.T) {
	out, err := FormatResults(sampleResults(), FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "key,status,env,value") {
		t.Errorf("CSV should start with header, got: %q", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in CSV output")
	}
}

func TestFormatResults_JSON(t *testing.T) {
	out, err := FormatResults(sampleResults(), FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "[") {
		t.Errorf("JSON should start with '[', got: %q", out)
	}
	if !strings.Contains(out, `"DB_HOST"`) {
		t.Errorf("expected DB_HOST key in JSON output")
	}
}

func TestFormatResults_JSONEmpty(t *testing.T) {
	out, err := FormatResults(nil, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "[]" {
		t.Errorf("expected empty JSON array, got: %q", out)
	}
}

func TestCSVEscape(t *testing.T) {
	cases := []struct{ in, want string }{
		{"simple", "simple"},
		{"with,comma", `"with,comma"`},
		{`has"quote`, `"has""quote"`},
	}
	for _, c := range cases {
		got := csvEscape(c.in)
		if got != c.want {
			t.Errorf("csvEscape(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
