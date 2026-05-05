package output

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch, Values: map[string]string{"prod": "localhost", "dev": "localhost"}},
		{Key: "API_KEY", Status: diff.StatusMissing, Values: map[string]string{"prod": "secret", "dev": ""}},
		{Key: "PORT", Status: diff.StatusMismatch, Values: map[string]string{"prod": "8080", "dev": "3000"}},
	}
}

func TestFormatResults_UnknownFormat(t *testing.T) {
	_, err := FormatResults(sampleResults(), "xml")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestFormatResults_EmptyText(t *testing.T) {
	out, err := FormatResults([]diff.Result{}, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestFormatResults_Text(t *testing.T) {
	out, err := FormatResults(sampleResults(), "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in text output")
	}
	if !strings.Contains(out, "<missing>") {
		t.Errorf("expected <missing> placeholder in text output")
	}
}

func TestFormatResults_CSV(t *testing.T) {
	out, err := FormatResults(sampleResults(), "csv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected header + data rows, got %d lines", len(lines))
	}
	if !strings.HasPrefix(lines[0], "key,status") {
		t.Errorf("expected CSV header to start with 'key,status', got %q", lines[0])
	}
}

func TestFormatResults_JSON(t *testing.T) {
	out, err := FormatResults(sampleResults(), "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"key"`) {
		t.Errorf("expected JSON key field, got %q", out)
	}
	if !strings.Contains(out, `"status"`) {
		t.Errorf("expected JSON status field, got %q", out)
	}
}

func TestFormatResults_Table(t *testing.T) {
	out, err := FormatResults(sampleResults(), "table")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected TABLE header KEY, got %q", out)
	}
	if !strings.Contains(out, "STATUS") {
		t.Errorf("expected TABLE header STATUS, got %q", out)
	}
}

func TestFormatResults_JSONEmpty(t *testing.T) {
	out, err := FormatResults([]diff.Result{}, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[]") {
		t.Errorf("expected empty JSON array, got %q", out)
	}
}
