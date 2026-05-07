package output

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestTableWriter_EmptyResults(t *testing.T) {
	// Use a strings.Builder directly instead.
	var sb strings.Builder
	w := NewTableWriter(&sb)
	w.Write([]diff.Result{})
	if !strings.Contains(sb.String(), "(no results)") {
		t.Errorf("expected '(no results)', got: %q", sb.String())
	}
}

func TestTableWriter_ShowsHeader(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch, Values: map[string]string{"prod": "localhost", "dev": "localhost"}},
	}
	var sb strings.Builder
	w := NewTableWriter(&sb)
	w.Write(results)
	out := sb.String()
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected header KEY, got: %q", out)
	}
	if !strings.Contains(out, "STATUS") {
		t.Errorf("expected header STATUS, got: %q", out)
	}
}

func TestTableWriter_ShowsRows(t *testing.T) {
	results := []diff.Result{
		{Key: "API_KEY", Status: diff.StatusMissing, Values: map[string]string{"prod": "secret", "dev": ""}},
		{Key: "PORT", Status: diff.StatusMismatch, Values: map[string]string{"prod": "8080", "dev": "3000"}},
	}
	var sb strings.Builder
	w := NewTableWriter(&sb)
	w.Write(results)
	out := sb.String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %q", out)
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %q", out)
	}
	if !strings.Contains(out, "missing") {
		t.Errorf("expected 'missing' status in output, got: %q", out)
	}
}

func TestTableWriter_EnvColumnsPresent(t *testing.T) {
	results := []diff.Result{
		{Key: "X", Status: diff.StatusMatch, Values: map[string]string{"alpha": "1", "beta": "1"}},
	}
	var sb strings.Builder
	w := NewTableWriter(&sb)
	w.Write(results)
	out := sb.String()
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected env name 'alpha' in header, got: %q", out)
	}
	if !strings.Contains(out, "beta") {
		t.Errorf("expected env name 'beta' in header, got: %q", out)
	}
}

func TestTableWriter_MismatchShowsValues(t *testing.T) {
	results := []diff.Result{
		{Key: "PORT", Status: diff.StatusMismatch, Values: map[string]string{"prod": "8080", "dev": "3000"}},
	}
	var sb strings.Builder
	w := NewTableWriter(&sb)
	w.Write(results)
	out := sb.String()
	if !strings.Contains(out, "8080") {
		t.Errorf("expected value '8080' in output, got: %q", out)
	}
	if !strings.Contains(out, "3000") {
		t.Errorf("expected value '3000' in output, got: %q", out)
	}
}
