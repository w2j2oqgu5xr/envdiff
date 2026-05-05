package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your/envdiff/internal/diff"
)

func makeSampleStats() []diff.EnvStats {
	return []diff.EnvStats{
		{EnvName: "prod", Total: 10, Matched: 9, Mismatched: 1, Missing: 0, Coverage: 100.0},
		{EnvName: "staging", Total: 10, Matched: 6, Mismatched: 2, Missing: 2, Coverage: 80.0},
		{EnvName: "dev", Total: 10, Matched: 3, Mismatched: 2, Missing: 5, Coverage: 50.0},
	}
}

func TestStatsPrinter_EmptyStats(t *testing.T) {
	var buf bytes.Buffer
	p := NewStatsPrinter(&buf)
	p.Print([]diff.EnvStats{})
	if !strings.Contains(buf.String(), "No environment statistics") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestStatsPrinter_ShowsHeader(t *testing.T) {
	var buf bytes.Buffer
	p := NewStatsPrinter(&buf)
	p.Print(makeSampleStats())
	out := buf.String()
	if !strings.Contains(out, "Environment Coverage Report") {
		t.Errorf("expected header in output, got: %s", out)
	}
}

func TestStatsPrinter_ShowsEnvNames(t *testing.T) {
	var buf bytes.Buffer
	p := NewStatsPrinter(&buf)
	p.Print(makeSampleStats())
	out := buf.String()
	for _, name := range []string{"prod", "staging", "dev"} {
		if !strings.Contains(out, name) {
			t.Errorf("expected env %q in output", name)
		}
	}
}

func TestStatsPrinter_ShowsCoverageValues(t *testing.T) {
	var buf bytes.Buffer
	p := NewStatsPrinter(&buf)
	p.Print(makeSampleStats())
	out := buf.String()
	for _, cov := range []string{"100.0%", "80.0%", "50.0%"} {
		if !strings.Contains(out, cov) {
			t.Errorf("expected coverage %q in output, got:\n%s", cov, out)
		}
	}
}

func TestStatsPrinter_NilWriterUsesStdout(t *testing.T) {
	// Just ensure it doesn't panic
	p := NewStatsPrinter(nil)
	if p == nil {
		t.Error("expected non-nil printer")
	}
}
