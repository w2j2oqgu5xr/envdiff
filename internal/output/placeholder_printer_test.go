package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeSamplePlaceholders() []diff.PlaceholderResult {
	return []diff.PlaceholderResult{
		{Key: "API_KEY", EnvName: "staging", Value: "CHANGEME", Pattern: "changeme"},
		{Key: "AUTH_TOKEN", EnvName: "development", Value: "<your-token>", Pattern: "<your"},
		{Key: "DB_PASS", EnvName: "staging", Value: "TODO", Pattern: "todo"},
	}
}

func TestPlaceholderPrinter_EmptyResults(t *testing.T) {
	buf := &bytes.Buffer{}
	p := &PlaceholderPrinter{w: buf}
	p.Print(nil)
	out := buf.String()
	if !strings.Contains(out, "No placeholder") {
		t.Errorf("expected no-placeholder message, got: %s", out)
	}
}

func TestPlaceholderPrinter_ShowsHeader(t *testing.T) {
	buf := &bytes.Buffer{}
	p := &PlaceholderPrinter{w: buf}
	p.Print(makeSamplePlaceholders())
	out := buf.String()
	if !strings.Contains(out, "Placeholder Values Detected") {
		t.Errorf("expected header in output, got: %s", out)
	}
}

func TestPlaceholderPrinter_ShowsKeys(t *testing.T) {
	buf := &bytes.Buffer{}
	p := &PlaceholderPrinter{w: buf}
	p.Print(makeSamplePlaceholders())
	out := buf.String()
	for _, key := range []string{"API_KEY", "AUTH_TOKEN", "DB_PASS"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %s in output", key)
		}
	}
}

func TestPlaceholderPrinter_ShowsEnvNames(t *testing.T) {
	buf := &bytes.Buffer{}
	p := &PlaceholderPrinter{w: buf}
	p.Print(makeSamplePlaceholders())
	out := buf.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("expected env name staging in output")
	}
	if !strings.Contains(out, "development") {
		t.Errorf("expected env name development in output")
	}
}

func TestPlaceholderPrinter_ShowsPatterns(t *testing.T) {
	buf := &bytes.Buffer{}
	p := &PlaceholderPrinter{w: buf}
	p.Print(makeSamplePlaceholders())
	out := buf.String()
	if !strings.Contains(out, "changeme") {
		t.Errorf("expected pattern changeme in output")
	}
	if !strings.Contains(out, "todo") {
		t.Errorf("expected pattern todo in output")
	}
}
