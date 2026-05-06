package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeSampleRedacts() []diff.RedactResult {
	return []diff.RedactResult{
		{Key: "API_KEY", EnvName: "prod", Original: "abc123", Redacted: "ab****"},
		{Key: "PASSWORD", EnvName: "dev", Original: "hunter2", Redacted: "hu*****"},
		{Key: "SECRET", EnvName: "prod", Original: "xyz", Redacted: "***"},
	}
}

func TestRedactPrinter_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	p := NewRedactPrinter(&buf)
	p.Print(nil)
	if !strings.Contains(buf.String(), "No sensitive keys") {
		t.Errorf("expected no-sensitive message, got: %q", buf.String())
	}
}

func TestRedactPrinter_ShowsHeader(t *testing.T) {
	var buf bytes.Buffer
	p := NewRedactPrinter(&buf)
	p.Print(makeSampleRedacts())
	if !strings.Contains(buf.String(), "Sensitive Keys Detected") {
		t.Errorf("expected header, got: %q", buf.String())
	}
}

func TestRedactPrinter_ShowsEnvNames(t *testing.T) {
	var buf bytes.Buffer
	p := NewRedactPrinter(&buf)
	p.Print(makeSampleRedacts())
	out := buf.String()
	if !strings.Contains(out, "prod") {
		t.Errorf("expected 'prod' in output")
	}
	if !strings.Contains(out, "dev") {
		t.Errorf("expected 'dev' in output")
	}
}

func TestRedactPrinter_ShowsRedactedValues(t *testing.T) {
	var buf bytes.Buffer
	p := NewRedactPrinter(&buf)
	p.Print(makeSampleRedacts())
	out := buf.String()
	if strings.Contains(out, "abc123") {
		t.Errorf("original value should not appear in output")
	}
	if !strings.Contains(out, "ab****") {
		t.Errorf("expected redacted value 'ab****' in output")
	}
}

func TestRedactPrinter_ShowsSummaryCount(t *testing.T) {
	var buf bytes.Buffer
	p := NewRedactPrinter(&buf)
	p.Print(makeSampleRedacts())
	out := buf.String()
	if !strings.Contains(out, "3") {
		t.Errorf("expected count '3' in summary output, got: %q", out)
	}
}

func TestRedactPrinter_NilWriter(t *testing.T) {
	p := NewRedactPrinter(nil)
	if p.w == nil {
		t.Error("expected non-nil writer after nil fallback")
	}
}
